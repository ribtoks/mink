# mink

[![Build Status](https://travis-ci.org/ribtoks/mink.svg?branch=master)](https://travis-ci.org/ribtoks/mink)
![license](https://img.shields.io/badge/license-MIT-blue.svg)
![copyright](https://img.shields.io/badge/%C2%A9-Taras_Kushnir-blue.svg)
![language](https://img.shields.io/badge/language-go-blue.svg)

## About

`mink` is a command line SEO tool that allows you to crawl URLs and get their basic metrics including, but not limited to: HTTP status code, meta description, size of the page, number of internal and external links and others.

It is a simple command-line alternative to tools like Screaming Frog SEO Spider, Netspeak Spider and other. It is useful to create plain-text or CSV report that can be used in spreadsheet software for further analysis.

## Install

`go get -u github.com/ribtoks/mink`

## Usage

```
Usage of mink:
  -d int
    	Maximum depth for crawling (default 1)
  -f string
    	Format of the output table|csv|tsv (default "table")
  -v	Write verbose logs
```

`mink` reads URLs from `STDIN` and writes reports to `STDOUT`. Report can be written in a form of a table, comma-separated values and tab-separated values.

## Examples

Crawl all pages of a single website:

`echo "https://your-website.com" | mink -d 1000 -f csv > report.csv`

Crawl a file with a list of URLs (1 per each line):

`cat urls.txt | mink -f csv > report.csv`

## Limitations

Currently mink does not handle javascript-based pages well.
