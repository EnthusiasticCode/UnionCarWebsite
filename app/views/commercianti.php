<article id="contacts" class="content-box">
	<h3>Area Commercianti</h3>
	<div class="row">
		<div class="small-12 large-6 columns">
			<h4>Hai già ricevuto l'autorizzazione come commerciante?</h4>
			<p>Se ti sei registrato come commerciante e la tua richiesta è stata approvata, dovresti aver ricevuto una e-mail con i dati di accesso da inserire qui sotto.</p>
			<form name="loginForm" class="custom" novalidate ng-submit="login.login()">
				<div class="row">
					<div class="small-8">
						<div class="row">
							<div class="small-3 columns">
								<label for="login-email" class="right inline">E-mail</label>
							</div>
							<div class="small-9 columns">
								<input required ng-model="login.email" type="email" name="email" id="login-email" ng-class="{ 'error':login.error }">
							</div>
						</div>
						<div class="row">
							<div class="small-3 columns">
								<label for="login-password" class="right inline">Password</label>
							</div>
							<div class="small-9 columns">
								<input required ng-model="login.password" type="password" name="password" id="login-password" ng-class="{ 'error':login.error }">
							</div>
						</div>
						<div class="row">
							<div class="small-12 columns">
								<input type="submit" value="Accedi" class="right login-submit small radius button" ng-class="{'disabled': !loginForm.$valid}">
							</div>
						</div>
					</div>
				</div>
			</form>
		</div>
		<div class="small-12 large-6 columns">
			<h4>Registrati come commerciante</h4>
			<form name="registerForm" class="custom" novalidate ng-submit="register.send()">
				<fieldset>
					<legend>Dati personali</legend>
					<div class="row">
						<div class="small-3 columns">
							<label for="register-name" class="right inline">Nome e cognome</label>
						</div>
						<div class="small-9 columns">
							<input id="register-name" required ng-model="register.name" type="text" name="name">
						</div>
					</div>
					<div class="row">
						<div class="small-3 columns">
							<label for="register-adderess" class="right inline">Indirizzo</label>
						</div>
						<div class="small-9 columns">
							<input id="register-adderess" required ng-model="register.adderess" type="text" name="adderess">
						</div>
					</div>
					<div class="row">
						<div class="small-3 columns">
							<label for="register-city" class="right inline">Città</label>
						</div>
						<div class="small-9 columns">
							<input id="register-city" required ng-model="register.city" type="text" name="city">
						</div>
					</div>
					<div class="row">
						<div class="small-3 columns">
							<label for="register-zip" class="right inline">CAP</label>
						</div>
						<div class="small-9 columns">
							<input id="register-zip" required ng-model="register.zip" type="text" name="zip">
						</div>
					</div>
					<div class="row">
						<div class="small-3 columns">
							<label for="register-region" class="right inline">Provincia</label>
						</div>
						<div class="small-9 columns">
							<input id="register-region" required ng-model="register.region" type="text" name="region">
						</div>
					</div>
					<div class="row">
						<div class="small-3 columns">
							<label for="register-vat" class="right inline">Partita Iva</label>
						</div>
						<div class="small-9 columns">
							<input id="register-vat" required ng-model="register.vat" type="text" name="vat">
						</div>
					</div>
				</fieldset>
				<fieldset>
					<legend>Dove inviare i dati d'accesso</legend>
					<div class="row">
						<div class="small-3 columns">
							<label for="register-email" class="right inline">E-mail</label>
						</div>
						<div class="small-9 columns">
							<input id="register-email" required ng-model="register.email" type="email" name="email">
						</div>
					</div>
					<div class="row">
						<div class="small-3 columns">
							<label for="register-password" class="right inline">Passowrd</label>
						</div>
						<div class="small-9 columns">
							<input id="register-password" required ng-model="register.password" type="password" name="password">
						</div>
					</div>
					<div class="row">
						<div class="small-3 columns">
							<label for="register-repassword" class="right inline">Ripeti password</label>
						</div>
						<div class="small-9 columns">
							<input id="register-repassword" required ng-model="register.repassword" type="password" name="repassword">
						</div>
					</div>
					<div class="row">
						<div class="small-12 columns">
							<small>
							<p>Ai sensi dell' Art.13 del D.lgs 196/03, la società UNION CAR SRL desidera informarLa che i Suoi dati personali forniti attraverso il presente sito Internet verranno acquisiti e trattati in forma cartacea e/o su supporto magnetico, elettronico o telematico nel pieno rispetto del Codice della Privacy.</p>

							<p>Il trattamento di tali dati potrà avvenire per finalità amministrative, gestionali, di selezione del personale, statistiche, commerciali e di marketing.</p>

							<p>Il conferimento dei dati stessi è pertanto facoltativo ed il suo rifiuto a fornirli e/o al successivo trattamento determinerà l'impossibilità per la scrivente di inserire i dati nel proprio archivio e conseguentemente instaurare eventuali rapporti con Lei.</p>

							<p>Relativamente ai dati medesimi Lei può esercitare i diritti previsti dall' art.7 del D.lgs n.196/2003 di cui per Sua opportuna informazione ne viene di seguito riportato il testo.</p>

							<p>Titolare del trattamento è UNION CAR SRL con sede in Via Cavin di Sala, 74 - 30035 Mirano - Venezia - Italy</p>

							<p>Ciò premesso, in mancanza di contrarie comunicazioni da parte Sua, consideriamo conferito a UNION CAR SRL il consenso all'utilizzo dei Suoi dati ai fini sopra indicati.</p>
							</small>
							<p><input type="checkbox" required ng-model="register.privacy" name="privacy"> Acconsento al trattamento dei dati personali</p>
						</div>
					</div>
				</fieldset>
				<div class="row">
					<div class="small-12 columns">
						<input type="submit" value="Registrati" class="right register-submit small radius button" ng-class="{'disabled': !registerForm.$valid}">
					</div>
				</div>
			</form>
		</div>
	</div>
</article>
