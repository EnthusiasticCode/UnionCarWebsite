<?php $car = $cars[0]; ?>

<article id="car-{{car.id}}" class="car-details">
	<header class="car-details-header">
		<a href="" class="secondary radius back button hide-for-small" onclick="return !window.history.back()"></a>
		<h1 class="title">
			<span class="brand" ng-bind="car.brand"><?php echo $car->brand; ?></span>
			/
			<strong class="model" ng-bind="car.model"><?php echo $car->model; ?></strong>
		</h1>
	</header>
	<div class="content-box">
		<div class="row">
			<div class="small-1 columns hide-for-small">
				<div class="brand-image brand-{{car.brand|lowercase}}"></div>
			</div>

			<div class="small-12 large-5 columns">
				<div class="description-container">
					<p ng-bind="car.model_description"></p>
				</div>
				<div class="price-container">
					<p>Prezzo al pubblico <span class="price">&euro; <strong ng-bind="car.price|number"><?php echo $car->price; ?></strong></span></p>
				</div>
				<ul class="details-list">
					<li>
						<div class="secondary radius label detail-label">Tipo</div>
						<div class="detail-value" ng-bind="car.car_type"><?php echo $car->car_type; ?></div>
					</li>
					<li>
						<div class="secondary radius label detail-label">Chilometraggio</div>
						<div class="detail-value"><span ng-bind="car.km|number"><?php echo $car->km; ?></span> Km</div>
					</li>
					<li>
						<div class="secondary radius label detail-label">Anno</div>
						<div class="detail-value" ng-bind="car.registration_date|date:'MM/yyyy'"><?php echo $car->registration_date; ?></div>
					</li>
					<li>
						<div class="secondary radius label detail-label">Potenza</div>
						<div class="detail-value"><span ng-bind="car.power_kw|number"><?php echo $car->power_kw; ?></span> KW (<span ng-bind="car.power_horses|number"><?php echo $car->power_horses; ?></span> CV)</div>
					</li>
				</ul>
				<h6>Ulteriori informazioni</h6>
				<ul class="details-list">
					<li>
						<div class="secondary radius label detail-label">Colore</div>
						<div class="detail-value" ng-bind="car.external_color|lowercase"><?php echo $car->external_color; ?></div>
					</li>
					<li>
						<div class="secondary radius label detail-label">Motore</div>
						<div class="detail-value"><span ng-bind="car.engine_size|number"><?php echo $car->engine_size; ?></span> cc</div>
					</li>
					<li>
						<div class="secondary radius label detail-label">Alimentazione</div>
						<div class="detail-value" ng-bind="car.fuel_type"><?php echo $car->fuel_type; ?></div>
					</li>
					<li>
						<div class="secondary radius label detail-label">Cambio</div>
						<div class="detail-value" ng-bind="car.gear_type"><?php echo $car->gear_type; ?></div>
					</li>
				</ul>
			</div>

			<div class="small-12 large-6 columns">
				<flex-slider slide="i in car.images" animation="slide" control-nav="thumbnails" ng-if="car.images.length">
					<li data-thumb="/car-images/{{i}}">
						<img ng-src="/car-images/{{i}}" alt="" />
					</li>
				</flex-slider>
				<?php for ($i=1; $i < 8; $i++) if ($car->{'image_url_'.$i}) : ?>
				<img src="/car-images/<?php echo $car->{'image_url_'.$i}; ?>" alt="" />
				<?php endif; ?>

				<a ng-href="/contacts/{{car.id}}" class="radius button contact-button">Contattaci per informazioni su <span class="brand" ng-bind="car.brand"><?php echo $car->brand; ?></span> <strong class="model" ng-bind="car.model"><?php echo $car->model; ?></strong></a>
			</div>
		</div> <!-- .row -->
	</div> <!-- .content-box -->
</article>
