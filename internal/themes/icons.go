package themes

import "fyne.io/fyne/v2"

// App logo - M with v in the valley
var resourceAppLogoSvg = &fyne.StaticResource{
	StaticName: "applogo.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
  <rect width="64" height="64" rx="12" fill="#1E88E5"/>
  <path d="M12 52V18h6l14 22 14-22h6v34h-7V30l-10.5 16h-5L19 30v22h-7z" fill="white"/>
  <path d="M32 10l-5 9h3l2-3.5 2 3.5h3l-5-9z" fill="white"/>
</svg>`),
}

// SVG icon resources for markdown editing toolbar

var resourceBoldSvg = &fyne.StaticResource{
	StaticName:    "bold.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M8 11h4.5a2.5 2.5 0 0 0 0-5H8v5Zm10 4.5a4.5 4.5 0 0 1-4.5 4.5H6V4h6.5a4.5 4.5 0 0 1 3.256 7.606A4.5 4.5 0 0 1 18 15.5ZM8 13v5h5.5a2.5 2.5 0 0 0 0-5H8Z"/></svg>`),
}

var resourceItalicSvg = &fyne.StaticResource{
	StaticName:    "italic.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M15 20H7v-2h2.927l2.116-12H9V4h8v2h-2.927l-2.116 12H15v2Z"/></svg>`),
}

var resourceHeadingSvg = &fyne.StaticResource{
	StaticName:    "heading.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M17 11V4h2v17h-2v-8H7v8H5V4h2v7h10Z"/></svg>`),
}

var resourceHeading1Svg = &fyne.StaticResource{
	StaticName:    "heading1.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M13 20h-2v-7H4v7H2V4h2v7h7V4h2v16zm8-12v12h-2v-9.796l-2 .536V8.67L19.5 8H21z"/></svg>`),
}

var resourceHeading2Svg = &fyne.StaticResource{
	StaticName:    "heading2.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M4 4v7h7V4h2v16h-2v-7H4v7H2V4h2zm14.5 4c2.071 0 3.75 1.679 3.75 3.75 0 .857-.288 1.648-.772 2.28l-.148.18L18.034 18H22v2h-7v-1.556l4.82-5.546c.268-.307.43-.709.43-1.148 0-.966-.784-1.75-1.75-1.75-.918 0-1.671.707-1.744 1.606l-.006.144h-2C14.75 9.679 16.429 8 18.5 8z"/></svg>`),
}

var resourceHeading3Svg = &fyne.StaticResource{
	StaticName:    "heading3.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M22 8l-.002 2-2.505 2.883c1.59.435 2.757 1.89 2.757 3.617 0 2.071-1.679 3.75-3.75 3.75-1.826 0-3.347-1.305-3.682-3.033l1.964-.382c.156.806.866 1.415 1.718 1.415.966 0 1.75-.784 1.75-1.75s-.784-1.75-1.75-1.75H17v-2h1.5c.966 0 1.75-.784 1.75-1.75s-.784-1.75-1.75-1.75c-.852 0-1.562.609-1.718 1.415l-1.964-.382C15.153 9.305 16.674 8 18.5 8c2.071 0 3.75 1.679 3.75 3.75 0 .475-.088.929-.249 1.347L22 13.097V8zM4 4v7h7V4h2v16h-2v-7H4v7H2V4h2z"/></svg>`),
}

var resourceLinkSvg = &fyne.StaticResource{
	StaticName:    "link.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M18.364 15.536 16.95 14.12l1.414-1.414a5 5 0 1 0-7.071-7.071L9.878 7.05 8.464 5.636 9.88 4.222a7 7 0 1 1 9.9 9.9l-1.415 1.414Zm-2.828 2.828-1.415 1.414a7 7 0 1 1-9.9-9.9l1.415-1.414L7.05 9.88l-1.414 1.414a5 5 0 1 0 7.071 7.071l1.414-1.414 1.415 1.414Zm-.708-10.607 1.415 1.415-7.071 7.07-1.415-1.414 7.071-7.07Z"/></svg>`),
}

var resourceCodeSvg = &fyne.StaticResource{
	StaticName:    "code.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="m23 12-7.071 7.071-1.414-1.414L20.172 12l-5.657-5.657 1.414-1.414L23 12ZM3.828 12l5.657 5.657-1.414 1.414L1 12l7.071-7.071 1.414 1.414L3.828 12Z"/></svg>`),
}

var resourceCodeBlockSvg = &fyne.StaticResource{
	StaticName:    "codeblock.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M3 3h18a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1zm1 2v14h16V5H4zm16 7-3.5 3.5-1.415-1.414L17.172 12l-2.087-2.086 1.415-1.414L20 12zm-11.414 0L6.5 14.086l-1.414-1.414L8.586 12l-3.5-3.5 1.414-1.414L10 10.586l-1.414 1.414zM11 18h2v-2h-2v2z"/></svg>`),
}

var resourceQuoteSvg = &fyne.StaticResource{
	StaticName:    "quote.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M4.583 17.321C3.553 16.227 3 15 3 13.011c0-3.5 2.457-6.637 6.03-8.188l.893 1.378c-3.335 1.804-3.987 4.145-4.247 5.621.537-.278 1.24-.375 1.929-.311 1.804.167 3.226 1.648 3.226 3.489a3.5 3.5 0 0 1-3.5 3.5 3.871 3.871 0 0 1-2.748-1.179Zm10 0C13.553 16.227 13 15 13 13.011c0-3.5 2.457-6.637 6.03-8.188l.893 1.378c-3.335 1.804-3.987 4.145-4.247 5.621.537-.278 1.24-.375 1.929-.311 1.804.167 3.226 1.648 3.226 3.489a3.5 3.5 0 0 1-3.5 3.5 3.871 3.871 0 0 1-2.748-1.179Z"/></svg>`),
}

var resourceHrSvg = &fyne.StaticResource{
	StaticName:    "hr.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M2 11h2v2H2v-2Zm4 0h12v2H6v-2Zm14 0h2v2h-2v-2Z"/></svg>`),
}

var resourceImageSvg = &fyne.StaticResource{
	StaticName:    "image.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M4.828 21l-.02.02-.021-.02H2.992A.993.993 0 0 1 2 20.007V3.993A1 1 0 0 1 2.992 3h18.016c.548 0 .992.445.992.993v16.014a1 1 0 0 1-.992.993H4.828zM20 15V5H4v14L14 9l6 6zm0 2.828l-6-6L6.828 19H20v-1.172zM8 11a2 2 0 1 1 0-4 2 2 0 0 1 0 4z"/></svg>`),
}

var resourceNewFileSvg = &fyne.StaticResource{
	StaticName:    "newfile.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M15 4H5v16h14V8h-4V4zM3 2.992C3 2.444 3.447 2 3.999 2H16l5 5v13.993A1 1 0 0 1 20.007 22H3.993A1 1 0 0 1 3 21.008V2.992zM11 11V8h2v3h3v2h-3v3h-2v-3H8v-2h3z"/></svg>`),
}

var resourceSaveAsSvg = &fyne.StaticResource{
	StaticName:    "saveas.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M18.172 7H6v6h12V7.828L16.172 6H6V4h11l3 3v12.993A1 1 0 0 1 18.993 21H5.007A1.001 1.001 0 0 1 4 19.993V4.007C4 3.451 4.449 3 5.007 3h11.586l1.707 1.707L20 6.414V8h-2V7h.172zM8 15v4h8v-4H8z"/></svg>`),
}

var resourceTableSvg = &fyne.StaticResource{
	StaticName:    "table.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M4 8h16V5H4v3zm10 11v-9h-4v9h4zm2 0h4v-9h-4v9zm-8 0v-9H4v9h4zM3 3h18a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z"/></svg>`),
}

var resourceLibrarySvg = &fyne.StaticResource{
	StaticName:    "library.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M3 18.5V5a3 3 0 0 1 3-3h14a1 1 0 0 1 1 1v18a1 1 0 0 1-1 1H6.5A3.5 3.5 0 0 1 3 18.5zM19 20v-3H6.5a1.5 1.5 0 0 0 0 3H19zM5 15.337A3.486 3.486 0 0 1 6.5 15H19V4H6a1 1 0 0 0-1 1v10.337z"/></svg>`),
}

var resourceSplitViewSvg = &fyne.StaticResource{
	StaticName:    "splitview.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M11 5H5v14h6V5zm2 0v14h6V5h-6zM4 3h16a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z"/></svg>`),
}

var resourceFocusSvg = &fyne.StaticResource{
	StaticName:    "focus.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M20 3H4a1 1 0 0 0-1 1v16a1 1 0 0 0 1 1h16a1 1 0 0 0 1-1V4a1 1 0 0 0-1-1zm-1 16H5V5h14v14z"/></svg>`),
}

var resourceHelpSvg = &fyne.StaticResource{
	StaticName:    "help.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M12 22C6.477 22 2 17.523 2 12S6.477 2 12 2s10 4.477 10 10-4.477 10-10 10zm0-2a8 8 0 1 0 0-16 8 8 0 0 0 0 16zm-1-5h2v2h-2v-2zm2-1.645V14h-2v-1.5a1 1 0 0 1 1-1 1.5 1.5 0 1 0-1.471-1.794l-1.962-.393A3.5 3.5 0 1 1 13 13.355z"/></svg>`),
}

var resourceSearchSvg = &fyne.StaticResource{
	StaticName:    "search.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M18.031 16.617l4.283 4.282-1.415 1.415-4.282-4.283A8.96 8.96 0 0 1 11 20c-4.968 0-9-4.032-9-9s4.032-9 9-9 9 4.032 9 9a8.96 8.96 0 0 1-1.969 5.617zm-2.006-.742A6.977 6.977 0 0 0 18 11c0-3.868-3.133-7-7-7-3.868 0-7 3.132-7 7 0 3.867 3.132 7 7 7a6.977 6.977 0 0 0 4.875-1.975l.15-.15z"/></svg>`),
}

var resourceStarSvg = &fyne.StaticResource{
	StaticName:    "star.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M12 18.26l-7.053 3.948 1.575-7.928L.587 8.792l8.027-.952L12 .5l3.386 7.34 8.027.952-5.935 5.488 1.575 7.928z"/></svg>`),
}

var resourcePresentationSvg = &fyne.StaticResource{
	StaticName:    "presentation.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M13 18v2h4v2H7v-2h4v-2H3a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h18a1 1 0 0 1 1 1v13a1 1 0 0 1-1 1h-8zM4 5v11h16V5H4z"/></svg>`),
}

var resourceSnippetSvg = &fyne.StaticResource{
	StaticName:    "snippet.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M16 2H8c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h8c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 18H8V4h8v16zM4 6H2v14c0 1.1.9 2 2 2h10v-2H4V6zm6 6h4v2h-4v-2zm0-4h4v2h-4V8z"/></svg>`),
}

var resourceSortSvg = &fyne.StaticResource{
	StaticName:    "sort.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M3 18h6v-2H3v2zM3 6v2h18V6H3zm0 7h12v-2H3v2z"/></svg>`),
}

var resourceGoalSvg = &fyne.StaticResource{
	StaticName:    "goal.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm-1-13h2v6h-2zm0 8h2v2h-2z"/></svg>`),
}

var resourceTypewriterSvg = &fyne.StaticResource{
	StaticName:    "typewriter.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M4 7h16v2H4V7zm0 4h16v2H4v-2zm0 4h10v2H4v-2zm16 0h-2v2h2v-2z"/></svg>`),
}

var resourcePrintSvg = &fyne.StaticResource{
	StaticName:    "print.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M19 8H5c-1.66 0-3 1.34-3 3v6h4v4h12v-4h4v-6c0-1.66-1.34-3-3-3zm-3 11H8v-5h8v5zm3-7c-.55 0-1-.45-1-1s.45-1 1-1 1 .45 1 1-.45 1-1 1zm-1-9H6v4h12V3z"/></svg>`),
}

var resourceExportSvg = &fyne.StaticResource{
	StaticName:    "export.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M19 12v7H5v-7H3v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2zm-6 .67l2.59-2.58L17 11.5l-5 5-5-5 1.41-1.41L11 12.67V3h2v9.67z"/></svg>`),
}
