/**
 * Pagination build function.  Malina eCommerce application
 *
 *
 * @author		John Aran (Ilyas Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */
package helper

import "fmt"

func PagingLinks(all, page, perPage int64, uri, attr,tag,class,classEl string) string {
	if all < 1 {
		return ""
	}
	var radius = int64(8) //app.Radius()

	var num int64
	if all % perPage == 0 {
		num = all / perPage
	} else {
		num = 1 + all /perPage
	}
	if num == 1 {
		return ""
	}

	var outStr = ""
	page = 1 + page /perPage
	if page < num+2 && page > 0 {
		var a int64

		var start int64
		if page > radius {
			start = page - radius
			if page >= radius+2 {
				//this.outstr += "1 href/"
				outStr += fmt.Sprintf(`<li `+classEl +`><`+tag +` `+class +` `+attr +`="`+uri +`"><<</`+tag +`></li>`, 0)
			}
		} else {
			start = 1
		}

		var fin int64
		if page +radius >= num {
			fin = num
		} else {
			fin = page + radius
		}
		if page > 1 {
			outStr += fmt.Sprintf(`<li `+classEl +`><`+tag +` `+class +` `+attr +`="`+uri +`"><</`+tag +`></li>`, (page -2)*perPage) //cur-1,
		}

		for i := start; i <= fin; i++ {

			if i == page {
				outStr += fmt.Sprintf(`<li class="active"><span>%d</span></li>`, i) //i
			} else {
				a = (i - 1) * perPage
				outStr += fmt.Sprintf(`<li `+classEl +`><`+tag +` `+class +` `+attr +`="`+uri +`">%d</`+tag +`></li>`, a, i)
			}
		}

		// End part
		if page < num {
			outStr += fmt.Sprintf(`<li `+classEl +`><`+tag +` `+class +` `+attr +`="`+uri +`">></`+tag +`></li>`, page * perPage) //cur+1
		}
		if num > page +radius {
			outStr += fmt.Sprintf(`<li `+classEl +`><`+tag +` `+class +` `+attr +`="`+uri +`">>> %d</`+tag +`></li>`, (num-1) * perPage, num) //)
		}
	}
	return outStr

}
