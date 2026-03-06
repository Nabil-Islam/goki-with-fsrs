package main

import (
	"regexp"
	"strings"
)

var symbolMap = map[string]string{
	// Greek lowercase
	`\alpha`:   "α", `\beta`: "β", `\gamma`: "γ", `\delta`: "δ",
	`\epsilon`: "ε", `\zeta`: "ζ", `\eta`:   "η", `\theta`: "θ",
	`\iota`:    "ι", `\kappa`: "κ", `\lambda`: "λ", `\mu`:   "μ",
	`\nu`:      "ν", `\xi`:   "ξ", `\pi`:    "π", `\rho`:   "ρ",
	`\sigma`:   "σ", `\tau`:  "τ", `\upsilon`: "υ", `\phi`:  "φ",
	`\chi`:     "χ", `\psi`:  "ψ", `\omega`: "ω",
	// Greek uppercase
	`\Gamma`: "Γ", `\Delta`: "Δ", `\Theta`: "Θ", `\Lambda`: "Λ",
	`\Xi`:    "Ξ", `\Pi`:    "Π", `\Sigma`: "Σ", `\Upsilon`: "Υ",
	`\Phi`:   "Φ", `\Psi`:   "Ψ", `\Omega`: "Ω",
	// Operators
	`\pm`: "±", `\times`: "×", `\div`: "÷", `\cdot`: "·",
	`\leq`: "≤", `\geq`: "≥", `\neq`: "≠", `\approx`: "≈",
	`\equiv`: "≡", `\sim`: "∼", `\propto`: "∝",
	// Sets
	`\in`: "∈", `\notin`: "∉", `\subset`: "⊂", `\supset`: "⊃",
	`\cup`: "∪", `\cap`: "∩", `\emptyset`: "∅",
	// Calculus/logic
	`\infty`: "∞", `\partial`: "∂", `\nabla`: "∇",
	`\int`: "∫", `\oint`: "∮", `\sum`: "Σ", `\prod`: "Π",
	`\forall`: "∀", `\exists`: "∃", `\neg`: "¬",
	`\land`: "∧", `\lor`: "∨", `\to`: "→", `\Rightarrow`: "⇒",
	`\Leftrightarrow`: "⟺", `\leftarrow`: "←", `\rightarrow`: "→",
	// Misc
	`\sqrt`: "√", `\therefore`: "∴", `\because`: "∵",
	`\ldots`: "…", `\cdots`: "⋯",
}

var (
	superMap = map[rune]string{
		'0': "⁰", '1': "¹", '2': "²", '3': "³", '4': "⁴",
		'5': "⁵", '6': "⁶", '7': "⁷", '8': "⁸", '9': "⁹",
		'n': "ⁿ", 'i': "ⁱ", '+': "⁺", '-': "⁻",
	}
	subMap = map[rune]string{
		'0': "₀", '1': "₁", '2': "₂", '3': "₃", '4': "₄",
		'5': "₅", '6': "₆", '7': "₇", '8': "₈", '9': "₉",
		'n': "ₙ", 'i': "ᵢ",
	}

	superRe = regexp.MustCompile(`\^\{([^}]+)\}|\^(\S)`)
	subRe   = regexp.MustCompile(`_\{([^}]+)\}|_(\S)`)
	fracRe  = regexp.MustCompile(`\\frac\{([^}]+)\}\{([^}]+)\}`)
)

func toSuper(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if u, ok := superMap[r]; ok {
			sb.WriteString(u)
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func toSub(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if u, ok := subMap[r]; ok {
			sb.WriteString(u)
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func UnicodeSubstitute(expr string) string {
	s := expr

	// simple fracs — a/b style
	s = fracRe.ReplaceAllStringFunc(s, func(m string) string {
		parts := fracRe.FindStringSubmatch(m)
		return parts[1] + "/" + parts[2]
	})

	// superscripts
	s = superRe.ReplaceAllStringFunc(s, func(m string) string {
		parts := superRe.FindStringSubmatch(m)
		content := parts[1]
		if content == "" {
			content = parts[2]
		}
		return toSuper(content)
	})

	// subscripts
	s = subRe.ReplaceAllStringFunc(s, func(m string) string {
		parts := subRe.FindStringSubmatch(m)
		content := parts[1]
		if content == "" {
			content = parts[2]
		}
		return toSub(content)
	})

	// symbol substitution — longest match first via sorted keys
	for sym, uni := range symbolMap {
		s = strings.ReplaceAll(s, sym, uni)
	}

	// clean up leftover braces
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")

	return s
}
