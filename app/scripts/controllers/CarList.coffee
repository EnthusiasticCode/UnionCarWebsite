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
			else if firstComponent() == "minPrice"
				minPrice = Number(otherComponents())
				if (template.minPrice or 0) < minPrice
					template.minPrice = minPrice
			else if firstComponent() == "maxPrice"
				maxPrice = Number(otherComponents())
				if (template.maxPrice or Infinity) > maxPrice
					template.maxPrice = maxPrice
			else if firstComponent() == "minKm"
				minKm = Number(otherComponents())
				if (template.minKm or 0) < minKm
					template.minKm = minKm
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
			text: 'Marca'
			children: [
				{ id: 'brand:Ferrari', text: 'Ferrari' }
				{ id: 'brand:BMW', text: 'BMW' } ]
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
