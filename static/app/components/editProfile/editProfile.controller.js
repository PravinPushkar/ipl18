'use strict';

var app = angular.module('ipl');

/**
 * Edit Profile Controller
 * 
 * Controller for the edit profile page.
 */
app.controller('editProfileController', function ($http, $mdToast, $scope, $window, urlService, utilsService, aliasPattern) {
    var vm = this;

    vm.edit = edit;
    vm.aliasPattern = aliasPattern;

    // Get the file name of image uploaded
    document.querySelector('input[id="profilePic"]').onchange = function () {
        vm.imageSelectedName = document.getElementById('profilePic').files[0].name;
    };

    // $scope.$on('fileProgress', function (e, progress) {
    //     vm.progress = progress.loaded / progress.total;
    // });

    // Send editable data to back-end
    function edit() {
        if ((vm.password !== '' && vm.password !== undefined && vm.password !== null) && (vm.alias !== '' && vm.alias !== undefined && vm.alias !== null) && !(document.getElementById('profilePic').files[0])) {
            utilsService.showToast({
                text: 'Please enter valid value in the fields.',
                hideDelay: 3000,
                isError: true
            });
            return;
        }
        if (vm.password !== vm.confirmPassword) {
            utilsService.showToast({
                text: 'Password and Confirm Password do not match',
                hideDelay: 3000,
                isError: true
            });
            return;
        }
        var fd = new FormData();
        if (vm.password !== '' && vm.password !== undefined && vm.password !== null) {
            fd.append('password', vm.password);
        }
        if (vm.alias !== '' && vm.alias !== undefined && vm.alias !== null) {
            fd.append('alias', vm.alias);
        }
        if (document.getElementById('profilePic').files[0]) {
            fd.append('image', document.getElementById('profilePic').files[0]);
        }
        var currentUserINumber = $window.localStorage.getItem('iNumber');
        var params = {
            url: `${urlService.userProfile}/${currentUserINumber}`,
            data: fd,
            method: 'PUT',
            transformRequest: angular.identity,
            headers: {
                'Content-Type': undefined
            }
        };
        $http(params)
            .then(function () {
                utilsService.showToast({
                    text: 'User Profile Updated.',
                    hideDelay: 3000,
                    isError: false
                });
                return;
            }, function (err) {
                utilsService.showToast({
                    text: `${err.message}.`,
                    hideDelay: 3000,
                    isError: true
                });
                return;
            });
    }
});