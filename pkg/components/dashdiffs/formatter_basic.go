package dashdiffs

import (
	"bytes"
	"html/template"
	diff "github.com/yudai/gojsondiff"
)

type BasicDiff struct {
	narrow		string
	keysIdent	int
	writing		bool
	LastIndent	int
	Block		*BasicBlock
	Change		*BasicChange
	Summary		*BasicSummary
}
type BasicBlock struct {
	Title		string
	Old			interface{}
	New			interface{}
	Change		ChangeType
	Changes		[]*BasicChange
	Summaries	[]*BasicSummary
	LineStart	int
	LineEnd		int
}
type BasicChange struct {
	Key			string
	Old			interface{}
	New			interface{}
	Change		ChangeType
	LineStart	int
	LineEnd		int
}
type BasicSummary struct {
	Key			string
	Change		ChangeType
	Count		int
	LineStart	int
	LineEnd		int
}
type BasicFormatter struct {
	jsonDiff	*JSONFormatter
	tpl			*template.Template
}

func NewBasicFormatter(left interface{}) *BasicFormatter {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tpl := template.Must(template.New("block").Funcs(tplFuncMap).Parse(tplBlock))
	tpl = template.Must(tpl.New("change").Funcs(tplFuncMap).Parse(tplChange))
	tpl = template.Must(tpl.New("summary").Funcs(tplFuncMap).Parse(tplSummary))
	return &BasicFormatter{jsonDiff: NewJSONFormatter(left), tpl: tpl}
}
func (b *BasicFormatter) Format(d diff.Diff) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_, err := b.jsonDiff.Format(d)
	if err != nil {
		return nil, err
	}
	bd := &BasicDiff{}
	blocks := bd.Basic(b.jsonDiff.Lines)
	buf := &bytes.Buffer{}
	err = b.tpl.ExecuteTemplate(buf, "block", blocks)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (b *BasicDiff) Basic(lines []*JSONLine) []*BasicBlock {
	_logClusterCodePath()
	defer _logClusterCodePath()
	blocks := make([]*BasicBlock, 0)
	for _, line := range lines {
		if b.returnToTopLevelKey(line) {
			if b.Block != nil {
				blocks = append(blocks, b.Block)
			}
		}
		b.LastIndent = line.Indent
		if line.Indent == 1 {
			if block, ok := b.handleTopLevelChange(line); ok {
				blocks = append(blocks, block)
			}
		}
		if line.Indent > 1 {
			if b.isSingleLineChange(line) {
				switch line.Change {
				case ChangeAdded, ChangeDeleted:
					b.Block.Changes = append(b.Block.Changes, &BasicChange{Key: line.Key, Change: line.Change, New: line.Val, LineStart: line.LineNum})
				case ChangeOld:
					b.Change = &BasicChange{Key: line.Key, Change: line.Change, Old: line.Val, LineStart: line.LineNum}
				case ChangeNew:
					b.Change.New = line.Val
					b.Change.LineEnd = line.LineNum
					b.Block.Changes = append(b.Block.Changes, b.Change)
				default:
				}
			} else {
				if line.Change != ChangeUnchanged {
					if line.Key != "" {
						b.narrow = line.Key
						b.keysIdent = line.Indent
					}
					if line.Change != ChangeNil {
						if !b.writing {
							b.writing = true
							key := b.Block.Title
							if b.narrow != "" {
								key = b.narrow
								if b.keysIdent > line.Indent {
									key = b.Block.Title
								}
							}
							b.Summary = &BasicSummary{Key: key, Change: line.Change, LineStart: line.LineNum}
						}
					}
				} else {
					if b.writing {
						b.writing = false
						b.Summary.LineEnd = line.LineNum
						b.Block.Summaries = append(b.Block.Summaries, b.Summary)
					}
				}
			}
		}
	}
	return blocks
}
func (b *BasicDiff) returnToTopLevelKey(line *JSONLine) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return b.LastIndent == 2 && line.Indent == 1 && line.Change == ChangeNil
}
func (b *BasicDiff) handleTopLevelChange(line *JSONLine) (*BasicBlock, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch line.Change {
	case ChangeNil:
		if line.Change == ChangeNil {
			if line.Key != "" {
				b.Block = &BasicBlock{Title: line.Key, Change: line.Change}
			}
		}
	case ChangeAdded, ChangeDeleted:
		return &BasicBlock{Title: line.Key, Change: line.Change, New: line.Val, LineStart: line.LineNum}, true
	case ChangeOld:
		b.Block = &BasicBlock{Title: line.Key, Old: line.Val, Change: line.Change, LineStart: line.LineNum}
	case ChangeNew:
		b.Block.New = line.Val
		b.Block.LineEnd = line.LineNum
		return b.Block, true
	default:
	}
	return nil, false
}
func (b *BasicDiff) isSingleLineChange(line *JSONLine) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return line.Key != "" && line.Val != nil && !b.writing
}

var (
	encStateMap	= map[ChangeType]string{ChangeAdded: "added", ChangeDeleted: "deleted", ChangeOld: "changed", ChangeNew: "changed"}
	tplFuncMap	= template.FuncMap{"getChange": func(c ChangeType) string {
		state, ok := encStateMap[c]
		if !ok {
			return "changed"
		}
		return state
	}}
)
var (
	tplBlock	= `{{ define "block" -}}
{{ range . }}
<div class="diff-group">
	<div class="diff-block">
		<h2 class="diff-block-title">
			<i class="diff-circle diff-circle-{{ getChange .Change }} fa fa-circle"></i>
			<strong class="diff-title">{{ .Title }}</strong> {{ getChange .Change }}
		</h2>


		<!-- Overview -->
		{{ if .Old }}
			<div class="diff-label">{{ .Old }}</div>
			<i class="diff-arrow fa fa-long-arrow-right"></i>
		{{ end }}
		{{ if .New }}
				<div class="diff-label">{{ .New }}</div>
		{{ end }}

		{{ if .LineStart }}
			<diff-link-json
				line-link="{{ .LineStart }}"
				line-display="{{ .LineStart }}{{ if .LineEnd }} - {{ .LineEnd }}{{ end }}"
				switch-view="ctrl.getDiff('html')"
			/>
		{{ end }}

	</div>

	<!-- Basic Changes -->
	{{ range .Changes }}
		<ul class="diff-change-container">
		{{ template "change" . }}
		</ul>
	{{ end }}

	<!-- Basic Summary -->
	{{ range .Summaries }}
		{{ template "summary" . }}
	{{ end }}

</div>
{{ end }}
{{ end }}`
	tplChange	= `{{ define "change" -}}
<li class="diff-change-group">
	<span class="bullet-position-container">
		<div class="diff-change-item diff-change-title">{{ getChange .Change }} {{ .Key }}</div>

		<div class="diff-change-item">
			{{ if .Old }}
				<div class="diff-label">{{ .Old }}</div>
				<i class="diff-arrow fa fa-long-arrow-right"></i>
			{{ end }}
			{{ if .New }}
					<div class="diff-label">{{ .New }}</div>
			{{ end }}
		</div>

		{{ if .LineStart }}
			<diff-link-json
				line-link="{{ .LineStart }}"
				line-display="{{ .LineStart }}{{ if .LineEnd }} - {{ .LineEnd }}{{ end }}"
				switch-view="ctrl.getDiff('json')"
			/>
		{{ end }}
	</span>
</li>
{{ end }}`
	tplSummary	= `{{ define "summary" -}}
<div class="diff-group-name">
	<i class="diff-circle diff-circle-{{ getChange .Change }} fa fa-circle-o diff-list-circle"></i>

	{{ if .Count }}
		<strong>{{ .Count }}</strong>
	{{ end }}

	{{ if .Key }}
		<strong class="diff-summary-key">{{ .Key }}</strong>
		{{ getChange .Change }}
	{{ end }}

	{{ if .LineStart }}
		<diff-link-json
			line-link="{{ .LineStart }}"
			line-display="{{ .LineStart }}{{ if .LineEnd }} - {{ .LineEnd }}{{ end }}"
			switch-view="ctrl.getDiff('json')"
		/>
	{{ end }}
</div>
{{ end }}`
)
