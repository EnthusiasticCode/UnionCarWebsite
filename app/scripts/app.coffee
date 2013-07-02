app = angular.module('UnionCarWebsiteApp', ['ngResource', 'angular-flexslider'])

app.config ($routeProvider, $locationProvider) ->
	$routeProvider.when '/',
		templateUrl: '/views/carlist.html'
		controller: 'CarListCtrl'
	.when '/details/:id',
		templateUrl: '/views/details.html',
		controller: 'DetailsCtrl'
	.when '/contacts',
		templateUrl: '/views/contacts.html',
		controller: 'ContactsCtrl'
	.otherwise
		redirectTo: '/'

	# see: http://ericduran.io/2013/05/31/angular-html5Mode-with-yeoman/
	# $locationProvider.html5Mode(yes).hashPrefix('!')
