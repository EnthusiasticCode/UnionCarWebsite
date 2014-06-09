<article id="contacts" class="content-box">
	<div class="row">
		<div class="small-12 large-6 columns">
			<h3>Union Car S.r.l</h3>
			<div class="row">
				<div class="small-12 large-5 columns">
					<a href="https://www.google.com/maps/preview#!q=Via+Cavin+di+Sala%2C+74%2C+Mirano%2C+30035+Province+of+Venice%2C+Italy&data=!4m10!1m9!4m8!1m3!1d340796!2d12.531552!3d45.467025!3m2!1i1440!2i730!4f13.1" target="_blank">
					<ul class="vcard address">
						<li class="street-address">Via Cavin di Sala 74</li>
						<li class="locality">Mirano</li>
						<li><span class="state">Venezia</span>, <span class="zip">30035</span></li>
					</ul>
					</a>
				</div>
				<div class="small-12 large-7 columns">
					<ul class="vcard contacts">
						<li><span class="secondary radius label">Telefono</span> <span>+39 041 570 3547</span></li>
						<li><span class="secondary radius label">Fax</span> <span>+39 041 570 3920</span></li>
						<li class="email"><span class="secondary radius label">E-Mail</span> <a href="mailto:mirano@unioncar.it">mirano@unioncar.it</a></li>
					</ul>
				</div>
			</div>
			<h4>Scrivici</h4>
			<div ng-switch="mail.status">
				<form name="contactForm" class="custom" novalidate ng-submit="mail.send()" ng-switch-default>
					<input type="email" placeholder="La tua email" name="email" ng-model="mail.data.sender" required ng-class="{'error': !contactForm.email.$error.required&&!contactForm.email.$valid}" class="contact-email">
					<small class="error" ng-show="!contactForm.email.$error.required&&!contactForm.email.$valid">E-Mail non valida</small>
					<textarea name="text" class="contact-text" placeholder="Messaggio" ng-model="mail.data.text" required></textarea>
					<input type="submit" value="Invia" class="contact-submit small radius button" ng-class="{'disabled': !contactForm.$valid}">
				</form>
				<div class="loading" ng-switch-when="sending"></div>
				<div ng-switch-when="sent">
					<p>Grazie per averci contattato!</p>
				</div>
				<div ng-switch-when="error">
					<p>Al momento non riusciamo ad inviare la tua richiesta!</p>
					<a href="" ng-href="mailto:mirano@unioncar.it?subject=Richiesta da {{mail.data.sender}}&amp;body={{mail.data.text}}" class="radius button">Inviala via E-Mail</a>
				</div>
			</div>
		</div>
		<div class="small-12 large-6 columns">
		</div>
	</div>
	<div id="map-canvas">
		<img src="/application/images/location.jpg" alt="">
	</div>
</article>
