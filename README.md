# mink

[![Build Status](https://travis-ci.org/ribtoks/mink.svg?branch=master)](https://travis-ci.org/ribtoks/mink)
![license](https://img.shields.io/badge/license-MIT-blue.svg)
![copyright](https://img.shields.io/badge/%C2%A9-Taras_Kushnir-blue.svg)
![language](https://img.shields.io/badge/language-go-blue.svg)

## About

`mink` is a command line SEO tool that allows you to crawl URLs and get their basic metrics including, but not limited to: HTTP status code, meta description, size of the page, number of internal and external links and others.

It is a simple command-line alternative to tools like Screaming Frog SEO Spider, Netspeak Spider and other.

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
