'use strict';

var app = angular.module('ipl');

/**
 * Themes
 * 
 * The theming config for the application
 */
app.config(['$mdThemingProvider', function ($mdThemingProvider) {
    // Add themes for toast
    $mdThemingProvider.theme('error-toast');
    $mdThemingProvider.theme('success-toast');
}]);