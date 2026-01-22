package gui

// ExportTheme represents an export theme style
type ExportTheme struct {
	Name        string
	Description string
	CSS         string
}

// GetExportThemes returns available export themes
func GetExportThemes() []ExportTheme {
	return []ExportTheme{
		{
			Name:        "Default",
			Description: "Clean, professional styling",
			CSS:         getDefaultExportCSS(),
		},
		{
			Name:        "GitHub",
			Description: "GitHub-flavored markdown style",
			CSS:         getGitHubExportCSS(),
		},
		{
			Name:        "Academic",
			Description: "Formal academic paper style",
			CSS:         getAcademicExportCSS(),
		},
		{
			Name:        "Dark",
			Description: "Dark theme for screen viewing",
			CSS:         getDarkExportCSS(),
		},
		{
			Name:        "Minimal",
			Description: "Clean, minimalist design",
			CSS:         getMinimalExportCSS(),
		},
		{
			Name:        "Print Friendly",
			Description: "Optimized for printing",
			CSS:         getPrintFriendlyCSS(),
		},
	}
}

func getDefaultExportCSS() string {
	return `body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    font-size: 16px;
    line-height: 1.6;
    max-width: 900px;
    margin: 0 auto;
    padding: 20px;
    color: #333;
}
h1, h2 { color: #2c5282; border-bottom: 1px solid #e2e8f0; padding-bottom: 0.3em; }
h3, h4 { color: #c05621; }
code {
    background-color: #f7fafc;
    padding: 0.2em 0.4em;
    border-radius: 3px;
    font-family: "SFMono-Regular", Consolas, monospace;
    font-size: 85%;
}
pre {
    background-color: #2d3748;
    color: #e2e8f0;
    padding: 16px;
    border-radius: 6px;
    overflow-x: auto;
}
pre code { background-color: transparent; padding: 0; color: inherit; }
blockquote {
    border-left: 4px solid #4299e1;
    margin: 0;
    padding-left: 16px;
    color: #4a5568;
    font-style: italic;
}
table { border-collapse: collapse; width: 100%; }
th, td { border: 1px solid #e2e8f0; padding: 8px 12px; text-align: left; }
th { background-color: #f7fafc; }
a { color: #4299e1; }
hr { border: none; border-top: 1px solid #e2e8f0; }`
}

func getGitHubExportCSS() string {
	return `body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    font-size: 16px;
    line-height: 1.5;
    max-width: 980px;
    margin: 0 auto;
    padding: 45px;
    color: #24292f;
}
h1, h2, h3, h4, h5, h6 {
    margin-top: 24px;
    margin-bottom: 16px;
    font-weight: 600;
    line-height: 1.25;
}
h1 { font-size: 2em; border-bottom: 1px solid #d0d7de; padding-bottom: 0.3em; }
h2 { font-size: 1.5em; border-bottom: 1px solid #d0d7de; padding-bottom: 0.3em; }
h3 { font-size: 1.25em; }
code {
    background-color: rgba(175,184,193,0.2);
    padding: 0.2em 0.4em;
    border-radius: 6px;
    font-family: ui-monospace, SFMono-Regular, monospace;
    font-size: 85%;
}
pre {
    background-color: #f6f8fa;
    padding: 16px;
    border-radius: 6px;
    overflow: auto;
    font-size: 85%;
    line-height: 1.45;
}
pre code { background-color: transparent; padding: 0; }
blockquote {
    border-left: 0.25em solid #d0d7de;
    color: #57606a;
    padding: 0 1em;
    margin: 0;
}
table { border-collapse: collapse; width: 100%; }
th, td { border: 1px solid #d0d7de; padding: 6px 13px; }
th { font-weight: 600; background-color: #f6f8fa; }
a { color: #0969da; text-decoration: none; }
a:hover { text-decoration: underline; }
hr { border: 0; border-top: 1px solid #d0d7de; margin: 24px 0; }
img { max-width: 100%; }`
}

func getAcademicExportCSS() string {
	return `body {
    font-family: "Times New Roman", Times, serif;
    font-size: 12pt;
    line-height: 2;
    max-width: 8.5in;
    margin: 1in auto;
    padding: 0;
    color: #000;
    text-align: justify;
}
h1 {
    font-size: 14pt;
    font-weight: bold;
    text-align: center;
    margin-top: 12pt;
    margin-bottom: 12pt;
}
h2 {
    font-size: 12pt;
    font-weight: bold;
    margin-top: 12pt;
    margin-bottom: 6pt;
}
h3 {
    font-size: 12pt;
    font-style: italic;
    margin-top: 12pt;
    margin-bottom: 6pt;
}
p { text-indent: 0.5in; margin: 0; }
p:first-of-type { text-indent: 0; }
code {
    font-family: "Courier New", monospace;
    font-size: 10pt;
}
pre {
    font-family: "Courier New", monospace;
    font-size: 10pt;
    margin: 12pt 0;
    padding: 12pt;
    background-color: #f5f5f5;
    border: 1px solid #ddd;
}
blockquote {
    margin: 12pt 0.5in;
    font-style: italic;
}
table {
    border-collapse: collapse;
    width: 100%;
    margin: 12pt 0;
}
th, td {
    border: 1px solid #000;
    padding: 6pt;
    text-align: left;
}
a { color: #000; }
hr { border: none; border-top: 1px solid #000; margin: 24pt 0; }
@page { margin: 1in; }`
}

func getDarkExportCSS() string {
	return `body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    font-size: 16px;
    line-height: 1.6;
    max-width: 900px;
    margin: 0 auto;
    padding: 40px;
    color: #c9d1d9;
    background-color: #0d1117;
}
h1, h2 { color: #58a6ff; border-bottom: 1px solid #30363d; padding-bottom: 0.3em; }
h3, h4 { color: #8b949e; }
code {
    background-color: rgba(110,118,129,0.4);
    padding: 0.2em 0.4em;
    border-radius: 6px;
    font-family: ui-monospace, SFMono-Regular, monospace;
    font-size: 85%;
}
pre {
    background-color: #161b22;
    padding: 16px;
    border-radius: 6px;
    overflow-x: auto;
}
pre code { background-color: transparent; padding: 0; }
blockquote {
    border-left: 4px solid #3b5998;
    margin: 0;
    padding-left: 16px;
    color: #8b949e;
}
table { border-collapse: collapse; width: 100%; }
th, td { border: 1px solid #30363d; padding: 8px 12px; }
th { background-color: #161b22; }
a { color: #58a6ff; }
hr { border: none; border-top: 1px solid #30363d; }
@media print {
    body { background-color: #fff; color: #000; }
    h1, h2 { color: #000; }
    pre { background-color: #f6f8fa; }
}`
}

func getMinimalExportCSS() string {
	return `body {
    font-family: Georgia, serif;
    font-size: 18px;
    line-height: 1.8;
    max-width: 700px;
    margin: 60px auto;
    padding: 20px;
    color: #333;
}
h1, h2, h3, h4, h5, h6 {
    font-family: -apple-system, sans-serif;
    font-weight: 600;
    margin-top: 2em;
    margin-bottom: 0.5em;
}
h1 { font-size: 2em; }
h2 { font-size: 1.5em; }
h3 { font-size: 1.2em; }
code {
    font-family: monospace;
    font-size: 0.9em;
    background: #f4f4f4;
    padding: 2px 6px;
}
pre {
    background: #f4f4f4;
    padding: 20px;
    overflow-x: auto;
    font-size: 0.9em;
}
pre code { background: none; padding: 0; }
blockquote {
    border-left: 3px solid #ccc;
    margin: 1.5em 0;
    padding-left: 1em;
    color: #666;
}
a { color: #0066cc; text-decoration: none; }
a:hover { text-decoration: underline; }
hr { border: none; border-top: 1px solid #eee; margin: 2em 0; }
img { max-width: 100%; height: auto; }`
}

func getPrintFriendlyCSS() string {
	return `body {
    font-family: Georgia, "Times New Roman", serif;
    font-size: 12pt;
    line-height: 1.5;
    max-width: 100%;
    margin: 0;
    padding: 0.5in;
    color: #000;
}
h1, h2, h3, h4, h5, h6 {
    font-family: Arial, Helvetica, sans-serif;
    page-break-after: avoid;
}
h1 { font-size: 18pt; margin-top: 0; }
h2 { font-size: 14pt; border-bottom: 1pt solid #000; }
h3 { font-size: 12pt; }
p { orphans: 3; widows: 3; }
code {
    font-family: "Courier New", monospace;
    font-size: 10pt;
    background: #eee;
    padding: 1pt 3pt;
}
pre {
    font-family: "Courier New", monospace;
    font-size: 9pt;
    background: #f5f5f5;
    border: 1pt solid #ccc;
    padding: 10pt;
    page-break-inside: avoid;
    white-space: pre-wrap;
    word-wrap: break-word;
}
pre code { background: none; padding: 0; }
blockquote {
    border-left: 2pt solid #666;
    margin: 0;
    padding-left: 10pt;
    color: #333;
}
table {
    border-collapse: collapse;
    width: 100%;
    page-break-inside: avoid;
}
th, td { border: 1pt solid #000; padding: 4pt 8pt; }
th { background: #eee; }
a { color: #000; text-decoration: underline; }
a[href]:after { content: " (" attr(href) ")"; font-size: 9pt; }
a[href^="#"]:after { content: ""; }
hr { border: none; border-top: 1pt solid #000; }
img { max-width: 100%; page-break-inside: avoid; }
@page { margin: 0.75in; }
@media screen {
    body { max-width: 800px; margin: 20px auto; }
}`
}
