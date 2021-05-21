package restful

import (
        "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"
        "github.com/go-openapi/loads"
        "github.com/go-openapi/runtime/middleware"
        "net/url"
        "os"
        "routing-manager/pkg/nbi"
        "routing-manager/pkg/restapi"
        "routing-manager/pkg/restapi/operations"
        "routing-manager/pkg/restapi/operations/debug"
        "routing-manager/pkg/restapi/operations/handle"
        "strconv"
)

func LaunchRest(nbiif string) {
    swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
    if err != nil {
        //log.Fatalln(err)
        xapp.Logger.Error(err.Error())
        os.Exit(1)
    }
    nbiUrl, err := url.Parse(nbiif)
    if err != nil {
        xapp.Logger.Error(err.Error())
        os.Exit(1)
    }
    api := operations.NewRoutingManagerAPI(swaggerSpec)
    server := restapi.NewServer(api)
    defer server.Shutdown()

    server.Port, err = strconv.Atoi(nbiUrl.Port())
    if err != nil {
        xapp.Logger.Error("Invalid NBI RestAPI port")
        os.Exit(1)
    }
    server.Host = "0.0.0.0"
   // set handlers
    api.HandleProvideXappHandleHandler = handle.ProvideXappHandleHandlerFunc(
        func(params handle.ProvideXappHandleParams) middleware.Responder {
            xapp.Logger.Info("Data received on Http interface")
            err := nbi.ProvideXappHandleHandlerImpl(params.XappCallbackData)
            if err != nil {
                xapp.Logger.Error("RoutingManager->AppManager request Failed: " + err.Error())
                return handle.NewProvideXappHandleBadRequest()
            } else {
                xapp.Logger.Info("RoutingManager->AppManager request Success")
                return handle.NewGetHandlesOK()
            }
        })
    api.HandleProvideXappSubscriptionHandleHandler = handle.ProvideXappSubscriptionHandleHandlerFunc(
        func(params handle.ProvideXappSubscriptionHandleParams) middleware.Responder {
            err := nbi.ProvideXappSubscriptionHandleImpl(params.XappSubscriptionData)
            if err != nil {
                xapp.Logger.Error("RoutingManager->SubManager Add Request Failed: " + err.Error())
                return handle.NewProvideXappSubscriptionHandleBadRequest()
            } else {
                xapp.Logger.Info("RoutingManager->SubManager Add Request Success, subid = %v, requestor = %v", *params.                  XappSubscriptionData.SubscriptionID, *params.XappSubscriptionData.Address)
                return handle.NewGetHandlesOK()
            }
        })
    api.HandleDeleteXappSubscriptionHandleHandler = handle.DeleteXappSubscriptionHandleHandlerFunc(
        func(params handle.DeleteXappSubscriptionHandleParams) middleware.Responder {
            err := nbi.DeleteXappSubscriptionHandleImpl(params.XappSubscriptionData)
            if err != nil {
                xapp.Logger.Error("RoutingManager->SubManager Delete Request Failed: " + err.Error())
                return handle.NewDeleteXappSubscriptionHandleNoContent()
           } else {
                xapp.Logger.Info("RoutingManager->SubManager Delete Request Success, subid = %v, requestor = %v", *params.               XappSubscriptionData.SubscriptionID, *params.XappSubscriptionData.Address)
                return handle.NewGetHandlesOK()
            }
        })
    api.HandleUpdateXappSubscriptionHandleHandler = handle.UpdateXappSubscriptionHandleHandlerFunc(
        func(params handle.UpdateXappSubscriptionHandleParams) middleware.Responder {
            err := nbi.UpdateXappSubscriptionHandleImpl(&params.XappList, params.SubscriptionID)
            if err != nil {
                return handle.NewUpdateXappSubscriptionHandleBadRequest()
            } else {
                return handle.NewUpdateXappSubscriptionHandleCreated()
            }
        })
    api.HandleCreateNewE2tHandleHandler = handle.CreateNewE2tHandleHandlerFunc(
        func(params handle.CreateNewE2tHandleParams) middleware.Responder {
            err := nbi.CreateNewE2tHandleHandlerImpl(params.E2tData)
            if err != nil {
                xapp.Logger.Error("RoutingManager->E2Manager AddE2T Request Failed: " + err.Error())
                return handle.NewCreateNewE2tHandleBadRequest()
            } else {
                xapp.Logger.Info("RoutingManager->E2Manager AddE2T Request Success, E2T = %v", *params.E2tData.E2TAddress)
                return handle.NewCreateNewE2tHandleCreated()
            }
        })

    api.HandleAssociateRanToE2tHandleHandler = handle.AssociateRanToE2tHandleHandlerFunc(
        func(params handle.AssociateRanToE2tHandleParams) middleware.Responder {
            err := nbi.AssociateRanToE2THandlerImpl(params.RanE2tList)
            if err != nil {
                xapp.Logger.Error("RoutingManager->E2Manager associateRanToE2T Request Failed: " + err.Error())
                return handle.NewAssociateRanToE2tHandleBadRequest()
            } else {
                xapp.Logger.Info("RoutingManager->E2Manager associateRanToE2T Request Success, E2T = %v", params.RanE2tList)
                return handle.NewAssociateRanToE2tHandleCreated()
            }
        })

    api.HandleDissociateRanHandler = handle.DissociateRanHandlerFunc(
        func(params handle.DissociateRanParams) middleware.Responder {
            err := nbi.DisassociateRanToE2THandlerImpl(params.DissociateList)
            if err != nil {
                xapp.Logger.Error("RoutingManager->E2Manager DisassociateRanToE2T Request Failed: " + err.Error())
                return handle.NewDissociateRanBadRequest()
            } else {
                xapp.Logger.Info("RoutingManager->E2Manager DisassociateRanToE2T Request Success, E2T = %v", params.DissociateList)
                return handle.NewDissociateRanCreated()
            }
        })
   api.HandleDeleteE2tHandleHandler = handle.DeleteE2tHandleHandlerFunc(
        func(params handle.DeleteE2tHandleParams) middleware.Responder {
            err := nbi.DeleteE2tHandleHandlerImpl(params.E2tData)
            if err != nil {
                xapp.Logger.Error("RoutingManager->E2Manager DeleteE2T Request Failed: " + err.Error())
                return handle.NewDeleteE2tHandleBadRequest()
            } else {
                xapp.Logger.Info("RoutingManager->E2Manager DeleteE2T Request Success, E2T = %v", *params.E2tData.E2TAddress)
                return handle.NewDeleteE2tHandleCreated()
            }
        })
    api.DebugGetDebuginfoHandler = debug.GetDebuginfoHandlerFunc(
        func(params debug.GetDebuginfoParams) middleware.Responder {
            response, err := nbi.DumpDebugData()
            if err != nil {
                return debug.NewGetDebuginfoCreated()
            } else {
                return debug.NewGetDebuginfoOK().WithPayload(&response)
            }
        })
    api.HandleAddRmrRouteHandler = handle.AddRmrRouteHandlerFunc(
        func(params handle.AddRmrRouteParams) middleware.Responder {
            err := nbi.Adddelrmrroute(params.RoutesList, true)
            if err != nil {
                return handle.NewAddRmrRouteBadRequest()
            } else {
                return handle.NewAddRmrRouteCreated()
            }

        })
    api.HandleDelRmrRouteHandler = handle.DelRmrRouteHandlerFunc(
        func(params handle.DelRmrRouteParams) middleware.Responder {
            err := nbi.Adddelrmrroute(params.RoutesList, false)
            if err != nil {
                return handle.NewDelRmrRouteBadRequest()
            } else {
                return handle.NewDelRmrRouteCreated()
            }
        })

    // start to serve API
    xapp.Logger.Info("Starting the HTTP Rest service")
    if err := server.Serve(); err != nil {
        xapp.Logger.Error(err.Error())
    }
}

