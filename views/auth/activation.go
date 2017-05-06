package authView


func ActivationForm(msg string)string{
	return `
	<!-- Activation form -->
	 <section>
		<div id="agileits-sign-in-page" class="sign-in-wrapper">
			<div class="agileinfo_signin">
				<h3>Activation</h3>

				<div class="alert alert-warning alert-dismissible" role="alert">
					<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
					<strong>`+msg+`</strong>
				</div>

			</div>
		</div>
	</section>
	<!-- //Activation form -->
	`
}



