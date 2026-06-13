package gui

import "regexp"

// mermaidBlockRe matches the HTML Goldmark emits for a ```mermaid fenced code
// block: <pre><code class="language-mermaid">…</code></pre>.
var mermaidBlockRe = regexp.MustCompile(`(?s)<pre><code class="language-mermaid">(.*?)</code></pre>`)

// convertMermaidBlocks rewrites Mermaid code blocks into the markup the
// Mermaid script renders: <pre class="mermaid">…</pre>. Without this, the
// injected mermaid.js (which scans for elements with class "mermaid") never
// finds the diagrams. The diagram source stays HTML-escaped; the browser
// decodes it when Mermaid reads the element's text content.
func convertMermaidBlocks(html string) string {
	return mermaidBlockRe.ReplaceAllString(html, `<pre class="mermaid">$1</pre>`)
}
