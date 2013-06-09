app = angular.module('UnionCarWebsiteApp', [])

app.config ['$routeProvider', ($routeProvider) ->
	$routeProvider.when '/',
		templateUrl: 'views/carlist.html'
		controller: 'CarListCtrl'
	.when '/details/:id',
		templateUrl: 'views/details.html',
		controller: 'DetailsCtrl'
	.otherwise
		redirectTo: '/']
