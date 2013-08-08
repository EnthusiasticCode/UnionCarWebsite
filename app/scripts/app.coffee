app = angular.module('UnionCarWebsiteApp', ['ngRoute', 'ngResource', 'angular-flexslider'])

app.config ($routeProvider, $locationProvider, ci) ->
	$routeProvider.when '/',
		templateUrl: "#{ci.appUrl}/views/carlist.php",
		controller: 'CarListCtrl'
	.when '/details/:id',
		templateUrl: "#{ci.appUrl}/views/details.php",
		controller: 'DetailsCtrl'
	.when '/contacts',
		templateUrl: "#{ci.appUrl}/views/contacts.php",
		controller: 'ContactsCtrl'
	.when '/contacts/:carId',
		templateUrl: "#{ci.appUrl}/views/contacts.php",
		controller: 'ContactsCtrl'
	.otherwise
		redirectTo: '/'

	$locationProvider.html5Mode(yes).hashPrefix('!')
