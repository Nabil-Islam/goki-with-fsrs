package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ---- Types ----

type SegmentType int

const (
	TextSegment SegmentType = iota
	InlineMath
	BlockMath
)

type Segment struct {
	Type    SegmentType
	Content string
}

// ---- LaTeX template ----

var latexTemplate = `\documentclass[preview,border=2pt]{standalone}
\usepackage{amsmath}
\usepackage{amssymb}
\begin{document}
$\displaystyle %s$
\end{document}`

// ---- Regexes ----

var blockRe  = regexp.MustCompile(`\$\$(.*?)\$\$`)
var inlineRe = regexp.MustCompile(`\$(.*?)\$`)

var complexPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\\frac`),
	regexp.MustCompile(`\\begin`),
	regexp.MustCompile(`\\matrix`),
	regexp.MustCompile(`\\sqrt\s*\[`),
	regexp.MustCompile(`\\int.*\\int`),
	regexp.MustCompile(`\{[^}]*\{[^}]*\}`),
	regexp.MustCompile(`\\(over|under)set`),
	regexp.MustCompile(`\\(lim|sum|prod)_\{[^}]+\}\^`),
}

// ---- Functions ----

func IsComplex(expr string) bool {
	for _, re := range complexPatterns {
		if re.MatchString(expr) {
			return true
		}
	}
	return false
}

func ParseSegments(content string) []Segment {
	var segments []Segment
	remaining := content

	for len(remaining) > 0 {
		loc := blockRe.FindStringIndex(remaining)
		inlineLoc := inlineRe.FindStringIndex(remaining)

		if loc != nil && (inlineLoc == nil || loc[0] <= inlineLoc[0]) {
			if loc[0] > 0 {
				segments = append(segments, Segment{TextSegment, remaining[:loc[0]]})
			}
			match := blockRe.FindStringSubmatch(remaining[loc[0]:loc[1]])
			segments = append(segments, Segment{BlockMath, match[1]})
			remaining = remaining[loc[1]:]
		} else if inlineLoc != nil {
			if inlineLoc[0] > 0 {
				segments = append(segments, Segment{TextSegment, remaining[:inlineLoc[0]]})
			}
			match := inlineRe.FindStringSubmatch(remaining[inlineLoc[0]:inlineLoc[1]])
			segments = append(segments, Segment{InlineMath, match[1]})
			remaining = remaining[inlineLoc[1]:]
		} else {
			segments = append(segments, Segment{TextSegment, remaining})
			break
		}
	}
	return segments
}

func RenderContent(content string) string {
	f, _ := os.OpenFile("/tmp/goki-debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	f.WriteString("RenderContent called: " + content + "\n")

	segments := ParseSegments(content)
	var sb strings.Builder

	for _, seg := range segments {
		f.WriteString(fmt.Sprintf("segment type=%d content=%s\n", seg.Type, seg.Content))
		switch seg.Type {
		case TextSegment:
			sb.WriteString(seg.Content)
		case InlineMath, BlockMath:
			f.WriteString("IsComplex: " + fmt.Sprint(IsComplex(seg.Content)) + "\n")
			if IsComplex(seg.Content) {
				sb.WriteString(UnicodeSubstitute(seg.Content))
				go OpenLatexPreview(seg.Content)
			} else {
				sb.WriteString(UnicodeSubstitute(seg.Content))
			}
		}
	}
	return sb.String()
}

func OpenLatexPreview(expr string) {
	f, _ := os.OpenFile("/tmp/goki-debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	log := func(s string) { f.WriteString(s + "\n") }
	log("OpenLatexPreview called with: " + expr)

	dir, err := os.MkdirTemp("", "goki-latex-*")
	if err != nil {
		log("MkdirTemp failed: " + err.Error())
		return
	}

	texFile := filepath.Join(dir, "eq.tex")
	pdfFile := filepath.Join(dir, "eq.pdf")
	pngFile := filepath.Join(dir, "eq.png")

	src := fmt.Sprintf(latexTemplate, expr)
	if err := os.WriteFile(texFile, []byte(src), 0644); err != nil {
		log("WriteFile failed: " + err.Error())
		return
	}

	cmd := exec.Command("pdflatex",
		"-interaction=nonstopmode",
		"-output-directory="+dir,
		texFile,
	)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log("pdflatex failed: " + err.Error() + "\n" + string(out))
		cmd = exec.Command("tectonic", "-o", dir, texFile)
		out, err = cmd.CombinedOutput()
		if err != nil {
			log("tectonic failed: " + err.Error() + "\n" + string(out))
			return
		}
	}

	cmd = exec.Command("convert", "-density", "200", pdfFile, pngFile)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log("convert failed: " + err.Error() + "\n" + string(out))
		exec.Command("open", pdfFile).Start()
		return
	}

	err = exec.Command("open", pngFile).Start()
	log("open result: " + fmt.Sprint(err) + " " + pngFile)
}
