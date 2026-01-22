package themes

import "fyne.io/fyne/v2"

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

var resourceLinkSvg = &fyne.StaticResource{
	StaticName:    "link.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="M18.364 15.536 16.95 14.12l1.414-1.414a5 5 0 1 0-7.071-7.071L9.878 7.05 8.464 5.636 9.88 4.222a7 7 0 1 1 9.9 9.9l-1.415 1.414Zm-2.828 2.828-1.415 1.414a7 7 0 1 1-9.9-9.9l1.415-1.414L7.05 9.88l-1.414 1.414a5 5 0 1 0 7.071 7.071l1.414-1.414 1.415 1.414Zm-.708-10.607 1.415 1.415-7.071 7.07-1.415-1.414 7.071-7.07Z"/></svg>`),
}

var resourceCodeSvg = &fyne.StaticResource{
	StaticName:    "code.svg",
	StaticContent: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor"><path d="m23 12-7.071 7.071-1.414-1.414L20.172 12l-5.657-5.657 1.414-1.414L23 12ZM3.828 12l5.657 5.657-1.414 1.414L1 12l7.071-7.071 1.414 1.414L3.828 12Z"/></svg>`),
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
