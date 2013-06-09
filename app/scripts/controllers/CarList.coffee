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

	$scope.search = {}
		# brand: ['Ferrari', 'bmw']
		# price: [25000, 1000000]
	$scope.carsFilterPredicate = (car) ->
		return no if $scope.search.brand? and $scope.search.brand.indexOf(car.brand) < 0
		return no if $scope.search.price? and not ($scope.search.price[0] <= car.price <= $scope.search.price[1])
		return no if $scope.search.km? and not ($scope.search.km[0] <= car.km <= $scope.search.km[1])
		yes

	$('.carlist-filter').select2
		tokenSeparators: [",", " "]
		escapeMarkup: (m) -> m
		formatResult: (item) ->
			# console.log item
			item.text
		formatSelection: (item) ->
			item.text
		# createSearchChoice: (item, data) ->
		# 	console.log arguments
		# 	return { id: term, text: term }
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