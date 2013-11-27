<div class="row">
	<div class="small-12 columns">
		<input type="text" id="carlist-filter" style="width:100%" data-placeholder="Filtra"/>
	</div>
</div>

<ul class="carlist content-box">

	<?php foreach($cars as $car): ?>

	<li ng-repeat="car in filteredCars = (cars|filter:carsFilterPredicate)" ng-animate="'animation'" id="car-{{car.id}}">
		<a href="/details/<?php echo $car->id; ?>" ng-href="/details/{{car.id}}" class="carlist-link">
			<div class="row">
				<div class="small-3 columns">
					<ul class="inline-list">
						<li class="carlist-brand-col">
							<div class="brand-image brand-{{car.brand|lowercase}}"></div>
						</li>
						<li class="carlist-image-col">
							<img
								ng-if="car.images.length"
								class="carlist-image"
								ng-src="/car-images/{{car.images[0]}}"
								<?php if ($car->image_url_1): ?>src="/car-images/<?php echo $car->image_url_1; ?>"<?php endif; ?>>
						</li>
					</ul>
				</div>
				<div class="small-9 columns">
					<div class="row">
						<h3 class="carlist-title">
							<span class="carlist-brand" ng-bind="car.brand"><?php echo $car->brand; ?></span>
							/
							<strong class="carlist-model" ng-bind="car.model"><?php echo $car->model; ?></strong>
						</h3>
					</div>
					<div class="row">
						<ul class="inline-list">
							<li class="carlist-type-col">
								<p class="carlist-type" ng-bind="car.car_type"><?php echo $car->car_type; ?></p>
							</li>
							<li class="carlist-date-col">
								<p class="carlist-date">anno <strong ng-bind="car.registration_date|date:'MM/yyyy'"><?php echo $car->registration_date; ?></strong></p>
							</li>
							<li class="carlist-price-col">
								<p class="carlist-price"><strong ng-bind="car.price&&('&euro; ' + (car.price|number))||'trattabile'"><?php echo $car->price; ?></strong></p>
							</li>
						</ul>
					</div>
				</div>
			</div>
		</a>
	</li>

	<?php endforeach; ?>

	<li ng-if="!filteredCars.length && filter.predicates.length" class="carlist-nothing-found">
		Nessuna auto trovata con i parametri indicati!
	</li>
	<li ng-if="!filteredCars.length && !filter.predicates.length" class="loading"></li>
</ul>

