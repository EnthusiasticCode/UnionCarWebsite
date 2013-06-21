app = angular.module('UnionCarWebsiteApp', [])

app.config ['$routeProvider', ($routeProvider) ->
	$routeProvider.when '/',
		templateUrl: 'views/carlist.html'
		controller: 'CarListCtrl'
	.when '/details/:id',
		templateUrl: 'views/details.html',
		controller: 'DetailsCtrl'
	.when '/contacts',
		templateUrl: 'views/contacts.html',
		controller: 'ContactsCtrl'
	.otherwise
		redirectTo: '/']
