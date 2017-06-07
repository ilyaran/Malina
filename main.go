/**
 * Main file for Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      il.aranov@gmail.com
 * @link
 * @github      https://github.com/ilyaran/Malina
 * @license	MIT License Copyright (c) 2017 John Aran (Ilyas Aranzhanovich Toxanbayev)
 */

package main

import (
	"fmt"
	"net/http"
	"github.com/ilyaran/Malina/controllers"
	"log"
	"github.com/gorilla/mux"
	"github.com/ilyaran/Malina/views/publicView"

	"github.com/ilyaran/Malina/core"
	"github.com/ilyaran/Malina/libraries"
)

func main() {
	fmt.Println("Listenning on port : 3001")

	core.MALINA = &core.Malina{}
	if core.MALINA.DB_init() {

		core.MALINA.LibrariesInit()
		core.MALINA.PublicViewObjectsInit()

		// uncomment if you need to generate crypted password string
		//generatePasswordForAnyAccount("ilyaran")

		router := mux.NewRouter().StrictSlash(true)
		router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/"))))

		router.HandleFunc("/", controller.PublicController.Index ).Methods("GET")

		router.HandleFunc("/public/product/{action:(?:list|get)}/", controller.PublicController.Index).Methods("GET")
		router.HandleFunc("/public/product/{action:(?:ajax)}/", controller.PublicController.Index).Methods("POST")
		router.HandleFunc("/public/cart/{action:(?:details)}/", controller.PublicController.Index).Methods("GET")
		router.HandleFunc("/public/cart/{action:(?:crud|ajax_list)}/", controller.PublicController.Index).Methods("POST")
		router.HandleFunc("/public/order/{action:(?:form)}/", controller.PublicController.Index).Methods("GET", "POST")

		router.HandleFunc("/auth/{action:(?:login|register|logout|forgot|change_password|activation|delete_account)}/", controller.AuthController.Index).Methods("POST","GET")

		router.HandleFunc("/home/cart/{action:(?:list|ajax_list|add|edit|del|get|inlist)}/", controller.Cart.Index).Methods("POST", "GET")
		router.HandleFunc("/home/product/{action:(?:list|ajax_list|add|edit|del|get|inlist)}/", controller.Product.Index).Methods("POST", "GET")
		router.HandleFunc("/home/category/{action:(?:list|ajax_list|add|edit|del|get|inlist)}/", controller.CategoryControllerObj.Index).Methods("POST", "GET")

		router.HandleFunc("/home/account/{action:(?:list|ajax_list|add|edit|del|get|inlist)}/", controller.AccountController.Index).Methods("POST", "GET")
		router.HandleFunc("/home/position/{action:(?:list)}/", controller.PositionController.Index).Methods("GET")
		router.HandleFunc("/home/position/{action:(?:list|ajax_list|inlist|add|edit|del|get)}/", controller.PositionController.Index).Methods("POST","GET")

		router.HandleFunc("/home/permission/{action:(?:list|ajax_list|inlist|add|edit|del|get)}/", controller.Permission.Index).Methods("POST","GET")


		router.HandleFunc("/filemanager/{action:(?:dirtree|createdir|deletedir|movedir|copydir|renamedir|fileslist|upload|download|downloaddir|deletefile|movefile|copyfile|renamefile)}/", controller.Filemanager.Index).Methods("POST")
		router.HandleFunc("/filemanager/{action:(?:thumb)}/", controller.Filemanager.Index).Methods("GET")

		log.Fatal(http.ListenAndServe(":3001", router))

		//core.MALINA.DB_close()
		fmt.Println("Exit")
	}
}
func generatePasswordForAnyAccount(pass string){
	fmt.Println("Your password: ",pass)
	fmt.Println("Crypted password: ",library.SESSION.Cryptcode(pass))
	fmt.Println(`Replace on database, table "account", column "account_password" with above crypted password`)
}
func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
			w.Write([]byte(publicView.Error(fmt.Sprintf("handling %q: %v", r.RequestURI, err))))
		}
	}
}

/*

.              any character, possibly including newline (flag s=true)
[xyz]          character class
[^xyz]         negated character class
\d             Perl character class
\D             negated Perl character class
[[:alpha:]]    ASCII character class
[[:^alpha:]]   negated ASCII character class
\pN            Unicode character class (one-letter name)
\p{Greek}      Unicode character class
\PN            negated Unicode character class (one-letter name)
\P{Greek}      negated Unicode character class
Composites:
xy             x followed by y
x|y            x or y (prefer x)
Repetitions:
x*             zero or more x, prefer more
x+             one or more x, prefer more
x?             zero or one x, prefer one
x{n,m}         n or n+1 or ... or m x, prefer more
x{n,}          n or more x, prefer more
x{n}           exactly n x
x*?            zero or more x, prefer fewer
x+?            one or more x, prefer fewer
x??            zero or one x, prefer zero
x{n,m}?        n or n+1 or ... or m x, prefer fewer
x{n,}?         n or more x, prefer fewer
x{n}?          exactly n x
Implementation restriction: The counting forms x{n,m}, x{n,}, and x{n} reject forms that create a minimum or maximum repetition count above 1000. Unlimited repetitions are not subject to this restriction.
Grouping:
(re)           numbered capturing group (submatch)
(?P<name>re)   named & numbered capturing group (submatch)
(?:re)         non-capturing group
(?flags)       set flags within current group; non-capturing
(?flags:re)    set flags during re; non-capturing

Flag syntax is xyz (set) or -xyz (clear) or xy-z (set xy, clear z). The flags are:

i              case-insensitive (default false)
m              multi-line mode: ^ and $ match begin/end line in addition to begin/end text (default false)
s              let . match \n (default false)
U              ungreedy: swap meaning of x* and x*?, x+ and x+?, etc (default false)
Empty strings:
^              at beginning of text or line (flag m=true)
$              at end of text (like \z not Perl's \Z) or line (flag m=true)
\A             at beginning of text
\b             at ASCII word boundary (\w on one side and \W, \A, or \z on the other)
\B             not at ASCII word boundary
\z             at end of text
Escape sequences:
\a             bell (== \007)
\f             form feed (== \014)
\t             horizontal tab (== \011)
\n             newline (== \012)
\r             carriage return (== \015)
\v             vertical tab character (== \013)
\*             literal *, for any punctuation character *
\123           octal character code (up to three digits)
\x7F           hex character code (exactly two digits)
\x{10FFFF}     hex character code
\Q...\E        literal text ... even if ... has punctuation
Character class elements:
x              single character
A-Z            character range (inclusive)
\d             Perl character class
[:foo:]        ASCII character class foo
\p{Foo}        Unicode character class Foo
\pF            Unicode character class F (one-letter name)
Named character classes as character class elements:
[\d]           digits (== \d)
[^\d]          not digits (== \D)
[\D]           not digits (== \D)
[^\D]          not not digits (== \d)
[[:name:]]     named ASCII class inside character class (== [:name:])
[^[:name:]]    named ASCII class inside negated character class (== [:^name:])
[\p{Name}]     named Unicode property inside character class (== \p{Name})
[^\p{Name}]    named Unicode property inside negated character class (== \P{Name})
Perl character classes (all ASCII-only):
\d             digits (== [0-9])
\D             not digits (== [^0-9])
\s             whitespace (== [\t\n\f\r ])
\S             not whitespace (== [^\t\n\f\r ])
\w             word characters (== [0-9A-Za-z_])
\W             not word characters (== [^0-9A-Za-z_])
ASCII character classes:
[[:alnum:]]    alphanumeric (== [0-9A-Za-z])
[[:alpha:]]    alphabetic (== [A-Za-z])
[[:ascii:]]    ASCII (== [\x00-\x7F])
[[:blank:]]    blank (== [\t ])
[[:cntrl:]]    control (== [\x00-\x1F\x7F])
[[:digit:]]    digits (== [0-9])
[[:graph:]]    graphical (== [!-~] == [A-Za-z0-9!"#$%&'()*+,\-./:;<=>?@[\\\]^_`{|}~])
[[:lower:]]    lower case (== [a-z])
[[:print:]]    printable (== [ -~] == [ [:graph:]])
[[:punct:]]    punctuation (== [!-/:-@[-`{-~])
[[:space:]]    whitespace (== [\t\n\v\f\r ])
[[:upper:]]    upper case (== [A-Z])
[[:word:]]     word characters (== [0-9A-Za-z_])
[[:xdigit:]]   hex digit (== [0-9A-Fa-f])
*/
/*
=====================================================================================
Email          string = "" // 1214 byte string literal not displayed
CreditCard     string = "" // 153 byte string literal not displayed
ISBN10         string = "^(?:[0-9]{9}X|[0-9]{10})$"
ISBN13         string = "^(?:[0-9]{13})$"
UUID3          string = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
UUID4          string = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
UUID5          string = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
UUID           string = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
Alpha          string = "^[a-zA-Z]+$"
Alphanumeric   string = "^[a-zA-Z0-9]+$"
Numeric        string = "^[0-9]+$"
Int            string = "^(?:[-+]?(?:0|[1-9][0-9]*))$"
Float          string = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
Hexadecimal    string = "^[0-9a-fA-F]+$"
Hexcolor       string = "^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$"
RGBcolor       string = "" // 159 byte string literal not displayed
ASCII          string = "^[\x00-\x7F]+$"
Multibyte      string = "[^\x00-\x7F]"
FullWidth      string = "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
HalfWidth      string = "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
Base64         string = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
PrintableASCII string = "^[\x20-\x7E]+$"
DataURI        string = "^data:.+\\/(.+);base64$"
Latitude       string = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
Longitude      string = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
DNSName        string = `^([a-zA-Z0-9]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9]{1}[a-zA-Z0-9_-]{1,62})*$`
IP             string = "" // 661 byte string literal not displayed
URLSchema      string = `((ftp|tcp|udp|wss?|https?):\/\/)`
URLUsername    string = `(\S+(:\S*)?@)`
Hostname       string = ``
URLPath        string = `((\/|\?|#)[^\s]*)`
URLPort        string = `(:(\d{1,5}))`
URLIP          string = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
URLSubdomain   string = `((www\.)|([a-zA-Z0-9]([-\.][-\._a-zA-Z0-9]+)*))`
URL            string = `^` + URLSchema + `?` + URLUsername + `?` + `((` + URLIP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + URLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + URLPort + `?` + URLPath + `?$`
SSN            string = `^\d{3}[- ]?\d{2}[- ]?\d{4}$`
WinPath        string = `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`
UnixPath       string = `^(/[^/\x00]*)+/?$`
Semver         string = "" // 185 byte string literal not displayed


*/
