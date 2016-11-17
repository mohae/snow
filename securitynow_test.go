package main

import (
	"bytes"
	"reflect"
	"testing"
)

var body = []byte(`
	<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
	<html xmlns="http://www.w3.org/1999/xhtml" lang="en" xml:lang="en">
	<head><title>GRC&nbsp;|&nbsp;Security Now! Episode Archive&nbsp;&nbsp;</title>
	<!-- ChangeDetection.com id=198mgjwxkml comment=" " detection="off" -->
	<meta name="author" content="Steve Gibson, GIBSON RESEARCH CORPORATION" />
	<meta name="keywords" content="Security Now!, weekly security audio column, steve gibson, internet security podcast, mp3 podcast, SockStress TCP Socket Stress and Exploitation, Denial of Service" />
	<meta name="description" content="Security Now! Weekly Internet Security Podcast. This week describing the newly revealed SockStress TCP stack vulnerabilities." />
	<meta name="sitemenu" content="" />
	<meta http-equiv="Content-Type" content="text/html; charset=ISO-8859-1" />
	<meta http-equiv="Content-Style-Type" content="text/css" />
	<meta http-equiv="pics-label" content='(pics-1.1 "http://www.rsac.org/ratingsv01.html" l gen true comment "RSACi North America Server" by "offices_@_grc.com" for "https://www.grc.com" on "1998.03.30T21:20-0800" r (n 0 s 0 v 0 l 0))' />
	<meta http-equiv="pics-label" content='(pics-1.1 "http://www.icra.org/ratingsv02.html" l gen true for "https://www.grc.com" r (cz 1 lz 1 nz 1 oz 1 vz 1) "http://www.rsac.org/ratingsv01.html" l gen true for "https://www.grc.com" r (n 0 s 0 v 0 l 0))' />
	<link rel="apple-touch-icon" sizes="76x76" href="/apple-touch-icon.png?v=m2d2o8jqMO">
	<link rel="icon" type="image/png" href="/favicon-32x32.png?v=m2d2o8jqMO" sizes="32x32">
	<link rel="icon" type="image/png" href="/favicon-16x16.png?v=m2d2o8jqMO" sizes="16x16">
	<link rel="manifest" href="/manifest.json?v=m2d2o8jqMO">
	<link rel="mask-icon" href="/safari-pinned-tab.svg?v=m2d2o8jqMO" color="#5bbad5">
	<link rel="shortcut icon" href="/favicon.ico?v=m2d2o8jqMO">
	<meta name="theme-color" content="#ffffff">
	<link rel="meta" href="labels.rdf" type="application/rdf+xml" title="ICRA labels" />
	<link rel="stylesheet" media="all" type="text/css" href="/grc.css" />
	<link rel="stylesheet" media="all" type="text/css" href="/mainmenu.css" />

	</head><body link="#CC0000" vlink="#006666" alink="#FF0000"><a name="top"></a>

	<!-- ########################## GRC Masthead Menu ########################## -->

	<div class="menuminwidth0"><div class="menuminwidth1"><div class="menuminwidth2">
	<div id="masthead">
		<a href="/default.htm"><img id="mastheadlogo" src="/mh-logo.gif" width="286" height="24" alt="Gibson Research Corporation" title="" /></a>
		<img id="focus" src="/mh-focus.gif" width="121" height="13" alt="" title="What we're about" />
		<a href="/news.htm"><img id="blogicon" src="/image/menublogicon.png" width="22" height="22" alt="blog icon" title="To our news and announcements page" /></a>
		<a href="/news.htm"><img id="twittericon" src="/image/menutwittericon.png" width="22" height="22" title="To our news and announcements page" alt="Twitter Icon"/></a>
		<a href="/news.htm"><img id="rssicon" src="/image/menurssicon.png" width="22" height="22" title="To our news and announcements page" alt="RSS Icon" /></a>

		<form action="https://www.google.com/search.htm" id="searchbox_000064552291181981813:y8yi5go2xza" onsubmit="javascript: this.action='https://www.grc.com/search.htm'">
			<input type="hidden" name="cx" value="000064552291181981813:y8yi5go2xza" />
			<input type="hidden" name="cof" value="FORID:11" />
			<input id="text" type="text" name="q" maxlength="256" />
			<input id="search" type="image" name="sa" value="Search" src="/mh-srch.gif" alt="&nbsp;[Search]" title="" />
		</form>
	</div>

	<div class="menu">

	<ul>
		<li><a href="/default.htm"><img src="/mb-home.gif" width="76" height="18" alt="[Home]" title="" /><!--[if gt IE 6]><!--></a><!--<![endif]--><!--[if lt IE 7]><table border="0" cellpadding="0" cellspacing="0"><tr><td><![endif]-->
			<ul class="leftbutton">
				<li><a href="/purchasing.htm">&nbsp;Purchasing</a></li>
				<li><a href="/sales.htm">&nbsp;Sales Support</a></li>
				<li><a href="/support.htm">&nbsp;Technical Support</a></li>
				<li><a href="/default.htm#bottom">&nbsp;Contact Us</a></li>
				<li><a href="/news.htm">&nbsp;Blogs, Twitter &amp; RSS</a></li>
				<li><a href="/privacy.htm">&nbsp;Privacy Policy</a></li>
	<!--			<li><a href="/siteoptions.htm">&nbsp;Site Options</a></li>			-->
				<li><a href="/stevegibson.htm">&nbsp;Steve's Projects Page</a></li>
				<li><a href="/resume.htm">&nbsp;Steve's Old Resume</a></li>
			</ul>
			<!--[if lte IE 6]></td></tr></table></a><![endif]-->
		</li>
	</ul>


	</td></tr></table></td></tr></table><img src="/image/transpixel.gif" width=1 height=8 border=0><br><font size=1>(Note that the text transcripts will appear a few hours later<br>than the audio files since they are created afterwards.)</font></center><br><b>For best results: RIGHT-CLICK on one of the two audio icons <img src="/image/speaker-hq.gif" width=16 height=16 hspace=2 border=0> &amp; <img src="/image/speaker-lq.gif" width=16 height=16 hspace=2 border=0> below</b> then choose "Save Target As..." to download the audio file to your computer <b>before</b> starting to listen. For the other resources you can either LEFT-CLICK to open in your browser or RIGHT-CLICK to save the resource to your computer.</font></td></tr></table>

	<!-- ################################################################################ -->

	</tr></table></td></tr></table></td></tr></table>

	<a name="500"></a>
	<br>
	<table width="85%" border=0 cellpadding=0 cellspacing=0><tr><td>
	<font size=1>Episode&nbsp;#500 | 24 Mar 2015 | 94 min.</font></td></tr></table><table width="85%" bgcolor="#999999" border=0 cellpadding=1 cellspacing=0><tr><td><table width="100%" bgcolor="#F8F8F8" border=0 cellpadding=0 cellspacing=0><tr><td><table width="100%" border=0 cellpadding=0 cellspacing=10><tr><td colspan=6><font size=1><font size=2><b>Windows Secure Boot</b></font><br><img src="/image/transpixel.gif" width=1 height=4 border=0><br>Leo and I discuss the recent Pwn2Own hacking competition.  We examine another serious breach of the Internet's certificate trust system and marvel at a very clever hack to crack the iPhone four-digit PIN lock.  Then we take a close look at the evolution of booting from BIOS to UEFI and how Microsoft has leveraged this into their &#8220;Windows Secure Boot&#8221; system.  We also examine what it might mean for the future of non-Windows operating systems.</font></td></tr><tr>

	<td><a href="https://media.grc.com/sn/sn-500.mp3"><img src="/image/speaker-hq.gif" width=16 height=16 title="RIGHT CLICK and SAVE AS to download a high quality MP3 audio file" border=0 align=left></a><font size=1>45&nbsp;MB</font></td>
	<td><a href="https://media.grc.com/sn/sn-500-lq.mp3"><img src="/image/speaker-lq.gif" width=16 height=16 title="RIGHT CLICK and SAVE AS to download a smaller MP3 audio file" border=0 align=left></a><font size=1>11&nbsp;MB</font></td>

	<td><a href="/sn/sn-500-notes.pdf"><img src="/image/snnotes.gif" width=16 height=16 title="LEFT CLICK on these to retrieve the PDF of Steve's show notes" border=0 align=left></a><font size=1>348&nbsp;KB</font></td>

	<td><a href="/sn/sn-500.htm"><img src="/image/htmlfile.gif" width=16 height=16 title="LEFT CLICK to view text transcript as web page" border=0 align=left></a><font size=1>126&nbsp;KB</font></td>
	<td><a href="/sn/sn-500.txt"><img src="/image/textfile.gif" width=16 height=16 title="RIGHT CLICK and SAVE AS to download a text format transcript" border=0 align=left></a><font size=1>73&nbsp;KB</font></td>
	<td><a href="/sn/sn-500.pdf"><img src="/image/pdffile.gif" width=16 height=16 title="RIGHT CLICK and SAVE AS to download PDF format transcript" border=0 align=left></a><font size=1>138&nbsp;KB</font></td>
	</tr></table></td></tr></table></td></tr></table>
`)

func TestGetLastEpisode(t *testing.T) {
	expected := 500
	expectedN := 7888
	var buf bytes.Buffer

	n, err := buf.Write(body)
	if err != nil {
		t.Error(err)
		return
	}
	if n != expectedN {
		t.Errorf("short write: wrote %d of %d", n, expectedN)
		return
	}
	tokens := getTokens(&buf)
	if len(tokens) == 0 {
		t.Error("got 0 tokens, expected some")
		return
	}
	n, err = lastEpisodeFromTokens(tokens)
	if err != nil {
		t.Error(err)
		return
	}
	if n != expected {
		t.Errorf("got %d; want %d", n, expected)
		return
	}
}

func TestSetEpisodeRange(t *testing.T) {
	tests := []struct {
		i           int
		cnf         Conf
		expected    Conf
		expectedErr string
	}{
		{
			i:           100,
			cnf:         Conf{lastN: 0, startEpisode: 0, stopEpisode: 0},
			expected:    Conf{lastN: 0, startEpisode: 1, stopEpisode: 100},
			expectedErr: "",
		},
		{
			i:           100,
			cnf:         Conf{lastN: 10, startEpisode: 0, stopEpisode: 0},
			expected:    Conf{lastN: 10, startEpisode: 91, stopEpisode: 100},
			expectedErr: "",
		},
		{
			i:           100,
			cnf:         Conf{lastN: 110, startEpisode: 0, stopEpisode: 0},
			expected:    Conf{lastN: 110, startEpisode: 1, stopEpisode: 100},
			expectedErr: "",
		},
		{
			i:           100,
			cnf:         Conf{lastN: 0, startEpisode: 110, stopEpisode: 0},
			expected:    Conf{lastN: 0, startEpisode: 0, stopEpisode: 0},
			expectedErr: "Nothing to do: the start episode, 110, does not yet exist. The last episode was 100.",
		},
		{
			i:           100,
			cnf:         Conf{lastN: 0, startEpisode: 42, stopEpisode: 0},
			expected:    Conf{lastN: 0, startEpisode: 42, stopEpisode: 100},
			expectedErr: "",
		},
		{
			i:           100,
			cnf:         Conf{lastN: 0, startEpisode: 42, stopEpisode: 101},
			expected:    Conf{lastN: 0, startEpisode: 42, stopEpisode: 100},
			expectedErr: "",
		},
		{
			i:           100,
			cnf:         Conf{lastN: 0, startEpisode: 11, stopEpisode: 42},
			expected:    Conf{lastN: 0, startEpisode: 11, stopEpisode: 42},
			expectedErr: "",
		},
		{
			i:           100,
			cnf:         Conf{lastN: 10, startEpisode: 11, stopEpisode: 42},
			expected:    Conf{lastN: 10, startEpisode: 11, stopEpisode: 42},
			expectedErr: "",
		},
	}

	for i, test := range tests {
		err := setEpisodeRange(test.i, &test.cnf)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%d: got %q; want %q", i, err.Error(), test.expectedErr)
			}
			continue
		}
		if !reflect.DeepEqual(test.cnf, test.expected) {
			t.Errorf("%d: got %v; want %v", i, test.cnf, test.expected)
		}
	}
}
