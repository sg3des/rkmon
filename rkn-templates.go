package main

import (
	"fmt"
	"io"
)

func head(w io.Writer) {
	w.Write(html0)
}

func (api *API) Template(w io.Writer) {
	head(w)
	w.Write(html1)
	fmt.Fprintf(w, `		<p>Blocked %v IP addresses</p>`, api.reg.TotalIP())
	fmt.Fprintf(w, `		<p>Last update at %v</p>`, api.CommitTime)
	w.Write(html2)
}

var html0 = []byte(`<!DOCTYPE html>
<html>
<head>
	<title>RKN monitor</title>
	<script src='/assets/index.js'></script>
	<link rel="stylesheet" type="text/css" href="/assets/index.css">
</head>
<body>`)
var html1 = []byte(`<main>
	<div>`)
var html2 = []byte(`	</div>
	<div id='status'></div>
	<form onsubmit="checkIP(this.ip.value, '#status'); return false;">
		<input type='text' name='ip' required>
		<input type='submit' value='submit'>
	</form>
</main>`)
