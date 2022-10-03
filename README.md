katch!
======
> a very simple wrapper utility for headless chrome to easily export any webpage as `png`, `pdf` or `html` (prerender), you can use it via `http` or `cli`.  `katch` was created after revamping [scraply](https://github.com/alash3al/scraply) to help in scraping client-side apps.

Usage
========

```bash
# run katch server in a separate terminal window
$ katch serve --listen=:3000
```

- a screenshot from google.com?
> `http://localhost:3000/export?format=png&url=https://google.com`

- a screenshot but with a custom screen size?
> `http://localhost:3000/export?format=png&viewport_width=1300&viewport_height=1600&url=https://google.com`  

- a full-screen screenshot!
> `http://localhost:3000/export?format=png&png_full_page=true&url=https://google.com`

- a pdf?
> `http://localhost:3000/export?format=pdf&pdf_print_background=true&pdf_landscape=false&pdf_paper_height=25&pdf_paper_width=20&viewport_width=1300&viewport_height=1600&url=https://google.com`

- as raw html?
> `http://localhost:3000/export?format=html&url=https://google.com`

- limit execution time to 5 seconds?
> `http://localhost:3000/export?format=pdf&max_exec_time=5&viewport_width=1300&viewport_height=1600&url=https://google.com`  

- scroll vertically by 100 pixels (scroll 1 time by 100 pixels vertically), you can do an infinity scroll by setting scroll_times=-1 !.
> `http://localhost:3000/export?format=png&url=https://youtube.com&scroll_step=100&scroll_times=1`

> **NOTE**: you can post JSON payload to the same endpoint with the same query params names above.

> **NOTE**: for `cli` usage, please run `$ katch help`


Download ?
==========
> you can go to the [releases page](https://github.com/alash3al/katch/releases) and pick the latest version.
> or you can `$ docker run --rm -it ghcr.io/alash3al/katch katch help`


Contribution ?
==============
> for sure you can contribute, how?

- clone the repo
- create your fix/feature branch
- create a pull request

nothing else, enjoy!

About
=====
> I'm [Mohamed Al Ashaal](https://alash3al.com), a software engineer :)
