package views

import (
	"Malina/config"
	"time"
	"fmt"
)

func Footer()string {
	return `
				</div>
			</div>
		</div>
	    	<script src="`+ app.Assets_backend_path() +`js/jquery.min.js"></script>

	<link rel="stylesheet" href="//code.jquery.com/ui/1.12.1/themes/base/jquery-ui.css">
  <script src="https://code.jquery.com/ui/1.12.1/jquery-ui.js"></script>
	    	<script src="`+ app.Assets_backend_path() +`js/bootstrap.min.js"></script>
	    	<script src="`+ app.Assets_backend_path() +`js/bootstrap-toggle.js"></script>
	    	<script src="`+app.Assets_path()+`js/home.js?`+fmt.Sprintf("%v",time.Now())+`"></script>
	</body>
</html>`
}