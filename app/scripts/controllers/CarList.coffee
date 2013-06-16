'use strict'

angular.module('UnionCarWebsiteApp')
.controller 'CarListCtrl', ['$scope', ($scope) ->
	$scope.cars = [
		{
			id: 1,
			brand: 'Ferrari',
			model: 'Testarossa',
			km: 110000,
			date: new Date(2003, 3),
			price: 230000,
			images: [
				'http://placehold.it/300&text=Ferrari'
			]
		},
		{
			id: 2,
			brand: 'Ferrari',
			model: 'Testarossa',
			km: 110000,
			date: new Date(2003, 3),
			price: 20000,
			images: [
				'http://placehold.it/300&text=Ferrari'
			]
		}
		]

	$scope.filter =
		template: {}
		predicates: []

	$scope.$watch("filter.predicates", (newValue, oldValue) ->
		template = {}
		angular.forEach(newValue, (predicateString) ->
			components = predicateString.split(":")
			firstComponent = ->
				components[0]
			otherComponents = ->
				components[1..].join(":")
			if components.length == 0
				return
			else if components.length == 1
				(template.texts or= []).push(firstComponent())
			else if firstComponent() == "text"
				(template.texts or= []).push(otherComponents())
			else if firstComponent() == "brand"
				(template.brands or= []).push(otherComponents())
			else if firstComponent() == "maxPrice"
				maxPrice = Number(otherComponents())
				if (template.maxPrice or Infinity) > maxPrice
					template.maxPrice = maxPrice
			else if firstComponent() == "maxKm"
				maxKm = Number(otherComponents())
				if (template.maxKm or Infinity) > maxKm
					template.maxKm = maxKm
			)
		$scope.filter.template = template
		)

	$scope.carsFilterPredicate = (car) ->
		template = $scope.filter.template
		if template.brands? and template.brands.indexOf(car.brand) < 0
			return no
		else if template.minPrice? and car.price < template.minPrice
			return no
		else if template.maxPrice? and car.price > template.maxPrice
			return no
		else if template.minKm? and car.km < template.minKm
			return no
		else if template.maxKm? and car.km > template.maxKm
			return no
		return yes

	$('.carlist-filter').select2(
		tokenSeparators: [",", " "]
		multiple: yes
		query: Select2.query.local([
			{
				text: 'Marca'
				children: [
					{ id: 'brand:Ferrari', text: 'Ferrari' }
					{ id: 'brand:BMW', text: 'BMW' }
				]
			},
			{
				text: 'Chilometraggio'
				children: [
					{ id: 'maxKm:2500', text: '2.500' }
					{ id: 'maxKm:5000', text: '5.000' }
					{ id: 'maxKm:10000', text: '10.000' }
					{ id: 'maxKm:20000', text: '20.000' }
					{ id: 'maxKm:30000', text: '30.000' }
					{ id: 'maxKm:40000', text: '40.000' }
					{ id: 'maxKm:50000', text: '50.000' }
					{ id: 'maxKm:60000', text: '60.000' }
					{ id: 'maxKm:70000', text: '70.000' }
					{ id: 'maxKm:80000', text: '80.000' }
					{ id: 'maxKm:90000', text: '90.000' }
					{ id: 'maxKm:100000', text: '100.000' }
					{ id: 'maxKm:125000', text: '125.000' }
					{ id: 'maxKm:150000', text: '150.000' }
					{ id: 'maxKm:175000', text: '175.000' }
					{ id: 'maxKm:200000', text: '200.000' }
				]
			},
			{
				text: 'Prezzo'
				children: [
					{ id: 'maxPrice:500', text: '500' }
					{ id: 'maxPrice:1000', text: '1.000' }
					{ id: 'maxPrice:1500', text: '1.500' }
					{ id: 'maxPrice:2000', text: '2.000' }
					{ id: 'maxPrice:2500', text: '2.500' }
					{ id: 'maxPrice:3000', text: '3.000' }
					{ id: 'maxPrice:4000', text: '4.000' }
					{ id: 'maxPrice:5000', text: '5.000' }
					{ id: 'maxPrice:6000', text: '6.000' }
					{ id: 'maxPrice:7000', text: '7.000' }
					{ id: 'maxPrice:8000', text: '8.000' }
					{ id: 'maxPrice:9000', text: '9.000' }
					{ id: 'maxPrice:10000', text: '10.000' }
					{ id: 'maxPrice:12500', text: '12.500' }
					{ id: 'maxPrice:15000', text: '15.000' }
					{ id: 'maxPrice:17500', text: '17.500' }
					{ id: 'maxPrice:20000', text: '20.000' }
					{ id: 'maxPrice:25000', text: '25.000' }
					{ id: 'maxPrice:30000', text: '30.000' }
					{ id: 'maxPrice:40000', text: '40.000' }
					{ id: 'maxPrice:50000', text: '50.000' }
					{ id: 'maxPrice:75000', text: '75.000' }
					{ id: 'maxPrice:100000', text: '100.000' }
				]
			}
		])
		createSearchChoice: (term) ->
			term = $.fn.select2.defaults.escapeMarkup(term)
			{ id: "text:#{term}", text: term }
		escapeMarkup: (m) -> m
		# formatResult: (item) ->
		# 	# console.log item
		# 	item.text
		# formatSelection: (item) ->
		# 	item.text
		).on('change', (e) -> $scope.$apply(-> $scope.filter.predicates = e.val))
	]

# Italian translation for select2
$.extend $.fn.select2.defaults,
	formatNoMatches: ->
		"Nessuna corrispondenza trovata"

	formatInputTooShort: (input, min) ->
		n = min - input.length
		"Inserisci ancora " + n + " caratter" + ((if n is 1 then "e" else "i"))

	formatInputTooLong: (input, max) ->
		n = input.length - max
		"Inserisci " + n + " caratter" + ((if n is 1 then "e" else "i")) + " in meno"

	formatSelectionTooBig: (limit) ->
		"Puoi selezionare solo " + limit + " element" + ((if limit is 1 then "o" else "i"))

	formatLoadMore: (pageNumber) ->
		"Caricamento in corso..."

	formatSearching: ->
		"Ricerca..."
