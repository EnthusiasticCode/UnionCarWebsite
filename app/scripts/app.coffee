app = angular.module('UnionCarWebsiteApp', [])

app.config ['$routeProvider', ($routeProvider) ->
	$routeProvider.when '/',
		templateUrl: 'views/carlist.html'
		controller: 'CarListCtrl'
	$routeProvider.otherwise
		redirectTo: '/']
