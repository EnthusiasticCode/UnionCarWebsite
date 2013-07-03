'use strict'

angular.module('UnionCarWebsiteApp')
.controller 'CarListCtrl', ($scope, carApi) ->
	$scope.cars = []
	carApi.all (cars) ->
		$scope.cars = cars

	$scope.filter =
		template: {}
		predicates: []
		select: []

	# Update filter template
	$scope.$watch "filter.predicates", (newValue, oldValue) ->
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
				(template.brands or= []).push(otherComponents().toLowerCase())
			else if firstComponent() == "maxPrice"
				maxPrice = Number(otherComponents())
				if (template.maxPrice or Infinity) > maxPrice
					template.maxPrice = maxPrice
			else if firstComponent() == "maxKm"
				maxKm = Number(otherComponents())
				if (template.maxKm or Infinity) > maxKm
					template.maxKm = maxKm
			else if firstComponent() == "minYear"
				minYear = Number(otherComponents())
				if (template.minYear or Infinity) > minYear
					template.minYear = minYear
			)
		$scope.filter.template = template

	# Predicate used by `|filter:` to filter cars based on filter.template
	$scope.carsFilterPredicate = (car) ->
		template = $scope.filter.template
		if template.brands? and template.brands.indexOf(car.brand?.toLowerCase()) < 0
			return no
		else if template.minPrice? and car.price < template.minPrice
			return no
		else if template.maxPrice? and car.price > template.maxPrice
			return no
		else if template.minKm? and car.km < template.minKm
			return no
		else if template.maxKm? and car.km > template.maxKm
			return no
		else if template.minYear? and car.date?.getFullYear() < template.minYear
			return no
		return yes

	# Update filter select from loaded cars
	$scope.$watch 'cars', (newCars) ->
		$scope.filter.select = [
			{
				text: 'Marca'
				children: ({ id: "brand:#{b}", text: b } for c in newCars when b = (c.brand?[0].toUpperCase()+c.brand?[1..-1].toLowerCase()))
			},
			{
				text: 'Prezzo fino a'
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
			},
			{
				text: 'Anno a partire da'
				children: [
					{ id: 'minYear:2013', text: '2013' }
					{ id: 'minYear:2012', text: '2012' }
					{ id: 'minYear:2011', text: '2011' }
					{ id: 'minYear:2010', text: '2010' }
					{ id: 'minYear:2009', text: '2009' }
					{ id: 'minYear:2008', text: '2008' }
					{ id: 'minYear:2007', text: '2007' }
					{ id: 'minYear:2006', text: '2006' }
					{ id: 'minYear:2005', text: '2005' }
					{ id: 'minYear:2004', text: '2004' }
					{ id: 'minYear:2003', text: '2003' }
					{ id: 'minYear:2002', text: '2002' }
					{ id: 'minYear:2001', text: '2001' }
					{ id: 'minYear:2000', text: '2000' }
					{ id: 'minYear:1999', text: '1999' }
					{ id: 'minYear:1998', text: '1998' }
					{ id: 'minYear:1997', text: '1997' }
					{ id: 'minYear:1996', text: '1996' }
					{ id: 'minYear:1995', text: '1995' }
					{ id: 'minYear:1994', text: '1994' }
					{ id: 'minYear:1993', text: '1993' }
					{ id: 'minYear:1992', text: '1992' }
					{ id: 'minYear:1991', text: '1991' }
					{ id: 'minYear:1990', text: '1990' }
					{ id: 'minYear:1989', text: '1989' }
					{ id: 'minYear:1988', text: '1988' }
					{ id: 'minYear:1987', text: '1987' }
					{ id: 'minYear:1986', text: '1986' }
					{ id: 'minYear:1985', text: '1985' }
					{ id: 'minYear:1980', text: '1980' }
					{ id: 'minYear:1975', text: '1975' }
					{ id: 'minYear:1970', text: '1970' }
					{ id: 'minYear:1965', text: '1965' }
					{ id: 'minYear:1960', text: '1960' }
					{ id: 'minYear:1955', text: '1955' }
					{ id: 'minYear:1950', text: '1950' }
					{ id: 'minYear:1940', text: '1940' }
					{ id: 'minYear:1930', text: '1930' }
					{ id: 'minYear:1920', text: '1920' }
					{ id: 'minYear:1910', text: '1910' }
				]
			}
		]

	# Update carlist select2 filter when filter.select chagnes
	$scope.$watch 'filter.select', (newSelect) ->
		$('#carlist-filter').select2('destroy').select2
			tokenSeparators: [",", " "]
			multiple: yes
			query: Select2.query.local(newSelect)
			createSearchChoice: (term) ->
				term = $.fn.select2.defaults.escapeMarkup(term)
				{ id: "text:#{term}", text: term }
			escapeMarkup: (m) -> m
			formatResult: (item, container, query, escapeMarkup) ->
				markup = []
				Select2.util.markMatch(item.text, query.term, markup, escapeMarkup)
				markup.unshift("&euro; ") if item.id?.indexOf('maxPrice') == 0
				markup.join("")
			formatSelection: (item) ->
				return "Fino a &euro; #{item.text}" if item.id.indexOf('maxPrice') == 0
				return "Prodotta dal #{item.text}" if item.id.indexOf('minYear') == 0
				item.text

	# Enable filtering change to update filter predicate
	$('#carlist-filter').on('change', (e) -> $scope.$apply(-> $scope.filter.predicates = e.val))

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
