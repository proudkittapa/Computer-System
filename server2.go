package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err.Error())
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	request(conn)
}

func request(conn net.Conn) {
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if i == 0 {
			mux(conn, ln)
		}
		if ln == "" {
			//headers are done
			break
		}
		i++
	}
}
func mux(conn net.Conn, ln string) {
	m := strings.Fields(ln)[0] //method
	u := strings.Fields(ln)[1] //url
	fmt.Println("***METHOD", m)
	fmt.Println("***URL", u)

	if m == "GET" && u == "/" {
		index(conn)
	}
}

func index(conn net.Conn) {
	body := `<!DOCTYPE html>
	<html lang="en">
	
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<!-- Bootstrap CSS -->
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css" integrity="sha384-JcKb8q3iqJ61gNV9KGb8thSsNjpSL0n8PARn9HuZOnIxN0hoP+VmmDGMN5t9UJ0Z" crossorigin="anonymous">
		<!-- title -->
		<title>Pin to Pre | Website for pre order</title>
		<link href="about_us.css" rel="stylesheet" type="text/css" />
	</head>
	
	
	<body>
	
		<div class=".container" >
			<!------------navigation bar(menu bar)------------>
			<div class="menu-bar">
				<div class="logo">
					<img src="images/PinToPre.png" width="125px">
		  </div>
				<nav>
				  <ul>
					<li>
					<a href="index.html">Home</a><i class="fa fa-home" aria-hidden="true"></i> </li>
					<li><a href="pre_order.html">Pre order</a><i class="fa fa-gratipay" aria-hidden="true"></i></li>
					<li><a href="pre_form.html">Promotion</a><i class="fa fa-tag" aria-hidden="true"></i></li>
					<li><a href="account.html">Account</a><i class="fa fa-user-circle" aria-hidden="true"></i></li>
					<li class="active" ><a href="about_us.html">About us</a><i class="fa fa-phone" aria-hidden="true"></i></li>
				  </ul>
				</nav>
			</div>
			<!------------about us paragraph------------>
			<div class="row">
				<div class="box-about">
				  <h1>About Us</h1><br>
				  <p>This text is styled with some of the text formatting properties. The heading uses the text-align, text-transform, and color properties.
				  The paragraph is indented, aligned, and the space between characters is specified. The underline is removed from this colored
				  <a target="_blank" href="tryit.asp?filename=trycss_text">"Try it Yourself"</a> link.</p>
				</div>
				<div class="header">
				<h1>need more information or anything u want to add eg.quote</h1>
				<p>pls write for me</p></div>
				<br><br>
			</div>
	  </div>
	
			<!------------Members------------>
	
		  <div class="container">
		   <div class="wrapper">
			  <h1>Team</h1>
				<div class="team">
					
					  <div class="card">
						<div class="team_img">
						  <img src="images/PinToPre.png" alt="team_member_img">
						</div>
						<h5>TitleTitle TitleTitleTitle</h5>
						<div>Subtitle Subtitle Subtitle Subtitle Subtitle Subtitle</div>
					  <div>Footer</div>
						
					  </div>
	
					
					  <div class="card">
						<div class="team_img">
						  <img src="images/PinToPre.png" alt="team_member_img">
						</div>
						<h5>Title</h5>
						<a>Subtitle</a>
						<div>Footer</div>
					  </div>
					
	
					  <div class="card">
						<div class="team_img">
						  <img src="images/PinToPre.png" alt="team_member_img">
						</div>
						<h5>Title</h5>
						<div>Subtitle</div>
						<div>Footer</div>
					  </div>
	
	
				</div>
			</div><br>
	</div>
		<!------------------------Thank you bottom------------------------>
	
			<div class="row">
				<div class="header">
					<h1>Thank you</h1>
					<p>for visiting our pages</p>
					<a href="pre_order.html" class="btn">more pre-order click!</a>
				</div>
				<br>
			</div>
	
		
	
	
	
	
	
		<!-- Optional JavaScript -->
		<!-- jQuery first, then Popper.js, then Bootstrap JS -->
		<script src="https://code.jquery.com/jquery-3.5.1.slim.min.js" integrity="sha384-DfXdz2htPH0lsSSs5nCTpuj/zy4C+OGpamoFVy38MVBnE+IbbVYUew+OrCXaRkfj" crossorigin="anonymous"></script>
		<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.1/dist/umd/popper.min.js" integrity="sha384-9/reFTGAW83EW2RDu2S0VKaIzap3H66lZH81PoYlFhbGU+6BZp6G7niu735Sk7lN" crossorigin="anonymous"></script>
		<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/js/bootstrap.min.js" integrity="sha384-B4gt1jrGC7Jh4AgTPSdUtOBvfO8shuf57BaghqFfPlYxofvL8/KUEfYiJOMMV+rV" crossorigin="anonymous"></script>
	
	  </body>
	</html>
	
	`
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprint(conn, body)
}
