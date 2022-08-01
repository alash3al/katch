katch!
======
> a very simple wrapper utility for headless chrome to easily export any webpage as `png`, `jpeg`, `pdf` or `html` (prerender), you can use it via `http` or `cli`.

Usage
========

```bash
# run katch server in a separate terminal window
$ katch serve --listen=:3000
```

- a screenshot from google.com?
> http://localhost:3000/export?format=png&url=https://google.com  

- a screenshot but with a big screen?
> http://localhost:3000/export?format=png&viewport_width=1300&viewport_height=1600&url=https://google.com  

- a pdf?
> http://localhost:3000/export?format=pdf&pdf_print_background=true&pdf_landscape=false&pdf_paper_height=25&pdf_paper_width=20&viewport_width=1300&viewport_height=1600&url=https://google.com

- limit execution time to 5 seconds?
> http://localhost:3000/export?format=pdf&max_exec_time=5&viewport_width=1300&viewport_height=1600&url=https://google.com  

> **NOTE**: you can post JSON payload to the same endpoint with the same name as the above query params.

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