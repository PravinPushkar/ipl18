'use strict';

var app = angular.module('ipl');

/**
 * Utils service
 * 
 * A collection of utility functions.
 */
app.factory('utilsService', ['$state', '$mdDialog', '$mdToast', '$window', function ($state, $mdDialog, $mdToast, $window) {
    var service = {};

    service.showConfirmDialog = showConfirmDialog;
    service.showToast = showToast;
    service.capitalizeFirstLetter = capitalizeFirstLetter;
    service.logout = logout;

    return service;

    // Function shows a confirm dialog
    function showConfirmDialog(params) {
        var confirm = $mdDialog.confirm()
            .title(params.title)
            .textContent(params.text)
            .ariaLabel(params.aria)
            .ok(params.ok)
            .cancel(params.cancel)
            .clickOutsideToClose(true);

        return $mdDialog.show(confirm);
    }

    // Function shows a toast
    function showToast(params) {
        $mdToast.show(
            $mdToast.simple()
            .position('fixed')
            .textContent(params.text)
            .hideDelay(params.hideDelay)
            .theme(params.isError ? 'error-toast' : 'success-toast')
            .action(params.hideDelay === 0 ? 'ok' : null)
        );
    }

    // Capitalizes first letter of an input string
    function capitalizeFirstLetter(str) {
        return str.charAt(0).toUpperCase() + str.slice(1);
    }

    // Logout function
    function logout(message, isError) {
        $state.go('login');
        $window.localStorage.removeItem('token');
        $window.localStorage.removeItem('iNumber');
        showToast({
            text: message,
            hideDelay: 1500,
            isError: isError
        });
    }
}]);