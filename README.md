Following the book [Hypermedia Systems](https://hypermedia.systems/) from CARSON GROSS, ADAM STEPINSKI and DENİZ AKŞİMŞEK.

Rebuilding the examples by using

- Go
- [htmx](https://htmx.org/)
- [Chi router](https://go-chi.io/#/)
- [templ](https://templ.guide/)
- SQLite
- [air](https://github.com/cosmtrek/air)

The branches of this project correspond to the chapters of the book.
The book's example backend is written in Python and Flask. This repo is an attempt to convert it to Go.

Branch for [Chapter 3](https://github.com/bigskysoftware/hypermedia-systems/blob/main/book/CH03_BuildingASimpleWebApplication.adoc) is an example of an old style (Web 1.0) contacts app, entirely server-side rendered using Go.
Deployed on Render: https://hypermedia-chapter3.onrender.com/

Branch for [Chapter 5](https://github.com/bigskysoftware/hypermedia-systems/blob/main/book/CH05_htmxPatterns.adoc) includes htmx patters like hx-boost for swapping the content of the html body only, hx-get/post/delete on issuing AJAX requests on click, change and keyup events, debouncing, selected swapping on a target, pagination, click to load and infinite scroll.
Deployed on Render: https://hypermedia-chapter5.onrender.com

Branch for [Chapter 6](https://github.com/bigskysoftware/hypermedia-systems/blob/main/book/CH06_MorehtmxPatterns.adoc) adds Active Search, single inline deletion of entries with swapping out the deleted content and bulk deletion.

Branch for [Chapter 7](https://github.com/bigskysoftware/hypermedia-systems/blob/main/book/CH07_ADynamicArchiveUIWithhtmx.adoc) added a feature do download the full contacts archive as JSON file. It included a UI component for triggering the archiving, polling for a (mocked) progress of the archiving and at last swapping the progress bar again for a download link.
Deployed on Render: https://hypermedia-chapter7.onrender.com

Branch for [Chapter 9](https://github.com/bigskysoftware/hypermedia-systems/blob/main/book/CH09_ScriptingInAHypermediaApplication.adoc) explored three lightweight scripting approaches - plain JavaScript, Alpine.js and \_hyperscript to add finishing touches to the UX.
Deployed on Render
