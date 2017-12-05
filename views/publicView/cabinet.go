/**
 * @author		John Aran (Ilyas Aranzhanovich Toxanbayev)
 * @version		1.0.0
 * @based on
 * @email      		il.aranov@gmail.com
 * @link
 * @github      	https://github.com/ilyaran/github.com/ilyaran/Malina
 * @license		MIT License Copyright (c) 2017 John Aran (Ilyas Toxanbayev)
 */

package publicView

import (
	"net/http"
	"github.com/ilyaran/Malina/berry"
	"html/template"
)

type Cabinet struct {

}


func (s *Cabinet)Index(malina *berry.Malina,w http.ResponseWriter){
	content :=`
	<div class="account_grid">
		<div class=" login-right">
			<div class="col-md-12">
				<h1>Cabinet</h1>


			</div>
		</div>
	</div>`
	malina.Content = template.HTML(content)
	PublicLayoutView.response(malina,w)
}