'use strict';

var app = angular.module('ipl');

/**
 * Utils service
 * 
 * A collection of utility functions.
 */
app.factory('utilsService', function ($state, $mdDialog, $mdToast) {
    var service = {};

    service.showConfirmDialog = showConfirmDialog;
    service.showToast = showToast;
    service.capitalizeFirstLetter = capitalizeFirstLetter;

    return service;

    // Function shows a confirm dialog
    function showConfirmDialog(params) {
        var confirm = $mdDialog.confirm()
            .title(params.title)
            .textContent(params.text)
            .ariaLabel(params.aria)
            .targetEvent(params.event)
            .ok(params.ok)
            .cancel(params.cancel)
            .clickOutsideToClose(true);

        return $mdDialog.show(confirm);
    }

    // Function shows a toast
    function showToast(params) {
        $mdToast.show(
            $mdToast.simple()
            .position('top right')
            .textContent(params.text)
            .hideDelay(params.hideDelay)
            .theme(params.isError ? 'error-toast' : 'success-toast')
            .action(params.hideDelay === 0 ? 'ok' : null)
        );
    }

    function capitalizeFirstLetter(str) {
        return str.charAt(0).toUpperCase() + str.slice(1);
    }
});