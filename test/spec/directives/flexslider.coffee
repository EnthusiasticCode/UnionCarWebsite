'use strict'

describe 'Directive: flexslider', () ->
  beforeEach module 'UnionCarWebsiteApp'

  element = {}

  it 'should make hidden element visible', inject ($rootScope, $compile) ->
    element = angular.element '<flexslider></flexslider>'
    element = $compile(element) $rootScope
    expect(element.text()).toBe 'this is the flexslider directive'
