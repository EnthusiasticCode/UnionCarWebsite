'use strict';

angular.module('UnionCarWebsiteApp')
	.directive 'flexSlider', ->
		restrict: 'E'
		scope: no
		replace: yes
		transclude: yes
		template: '<div class="flexslider-container"></div>'
		compile: (element, attr, linker) -> ($scope, $element, $attr) ->
			match = $attr.slide.match /^\s*(.+)\s+in\s+(.*?)\s*(\s+track\s+by\s+(.+)\s*)?$/
			indexString = match[1]
			collectionString = match[2]
			elementsScopes = []
			flexsliderDiv = null

			$scope.$watchCollection collectionString, (collection) ->
				# Remove old flexslider
				if elementsScopes.length > 0 or flexsliderDiv?
					$element.children().remove()
					for e in elementsScopes
						e.$destroy()
					elementsScopes = []

				# Create flexslider container
				slides = $('<ul class="slides"></ul>')
				flexsliderDiv = $('<div class="flexslider"></div>')
				flexsliderDiv.append slides
				$element.append flexsliderDiv

				# Early exit if no collection
				return unless collection?

				# Generate slides
				for c in collection
					childScope = $scope.$new()
					childScope[indexString] = c
					linker childScope, (clone) ->
						slides.append clone
						elementsScopes.push childScope

				# Running flexslider
				# Options are derived from flex-slider arguments
				setTimeout (-> $scope.$apply -> flexsliderDiv.flexslider $attr), 0
