..
..  Copyright (c) 2019 AT&T Intellectual Property.
..  Copyright (c) 2019 Nokia.
..
..  Licensed under the Creative Commons Attribution 4.0 International
..  Public License (the "License"); you may not use this file except
..  in compliance with the License. You may obtain a copy of the License at
..
..    https://creativecommons.org/licenses/by/4.0/
..
..  Unless required by applicable law or agreed to in writing, documentation
..  distributed under the License is distributed on an "AS IS" BASIS,
..  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
..
..  See the License for the specific language governing permissions and
..  limitations under the License.
..


Installation Guide
==================

.. contents::
   :depth: 3
   :local:

Introduction
------------
Routing Manager is a basic platform service of RIC. It is responsible for distributing routing policies to other platform components and xApps.

Installation of routing manager in a k8s environment
====================================================

This deployment guide is for testing routing manager in actual k8s environment.

.. include:: ./installation-guide-rtmgr-as-pod.rst


Installation of routing manager in a standalone linux machine
=============================================================

This deployment guide will be useful for development purposes

.. include:: ./installation-guide-rtmgr-as-exe.rst

