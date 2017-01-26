(function(root, factory) {
    if ("function" === typeof define && define.amd) {
        define("remote_manager",
            [
                "google_maps",
                "tile_system",
                "reconnecting-websocket",
                "inherence"
            ],
            factory
        );
    } else if ("object" === typeof module && module.exports) {
        module.exports = factory(
            require('google_maps'),
            require('tile_system'),
            require('reconnecting-websocket'),
            require('inherence')
        );
    } else {
        root.RemoteManager = factory(
            root.google,
            root.TileSystem,
            root.ReconnectingWebSocket,
            root.Inherence
        );
    }
})(this, function(google, TileSystem, ReconnectingWebSocket, Inherence) {
    var RemoteManager = function(map, wsPath, onPointReceived) {
        this.map = map;
        this.ws = new ReconnectingWebSocket(wsPath);
        this.onPointReceived = onPointReceived;
    };

    RemoteManager.prototype.run = function() {
        this.map.addListener('idle', this._onIdle.bind(this));
        this.ws.onmessage = this._onMessage.bind(this);
        this.ws.onclose = (function() {
            this.ws = new WebSocket(this.ws.url);
        }).bind(this);
    };

    RemoteManager.prototype._onIdle = function() {
        this._updateState();
        this._fetchPoints();
    };

    RemoteManager.prototype._fetchPoints = function() {
        var bounds = this.bounds, zoom = this.zoom;
        this._send({
            lat1: bounds.getSouthWest().lat(),
            lng1: bounds.getSouthWest().lng(),
            lat2: bounds.getNorthEast().lat(),
            lng2: bounds.getNorthEast().lng(),
            zoom: zoom
        });
    };

    RemoteManager.prototype._send = function(toSend) {
        this.ws.send(JSON.stringify(toSend));
    };

    RemoteManager.prototype._onMessage = function(event) {
        var data = JSON.parse(event.data);
        this.onPointReceived(data);
    };

    RemoteManager.prototype._updateState = function() {
        this.zoom = this.map.getZoom();
        this.bounds = this.map.getBounds();
    };

    var CachedRemoteManager = (function() {
        Inherence(CachedRemoteManager, RemoteManager);

        function CachedRemoteManager() {
            CachedRemoteManager.__super__.constructor.apply(this, arguments);
            var that = this;
            this.cache = {};
            this.originalCallback = this.onPointReceived;
            this.onPointReceived = (function(postHook, func) {
                return function(data, ignoreHook) {
                    func.call(this, data);
                    if (!ignoreHook) {
                        postHook.call(this, data);
                    }
                };
            })(this._postHookOnPointReceived, this.onPointReceived);
        }

        return CachedRemoteManager;
    })();

    CachedRemoteManager.prototype._updateState = function() {
        var zoom = this.zoom;
        CachedRemoteManager.__super__._updateState.apply(this);
        if (zoom && zoom != this.zoom) {
            this.cache = {};
        }
    };

    CachedRemoteManager.prototype._postHookOnPointReceived = function(data) {
        var zoom = this.zoom, tile, key;
        for (var k in data) {
            key = TileSystem.getQuadKey(data[k].lat, data[k].lng, zoom);
            if (!this.cache.hasOwnProperty(key)) {
                this.cache[key] = [];
            }

            this.cache[key].push(data[k]);
        }
    };

    CachedRemoteManager.prototype._fetchPoints = function() {
        var bounds = this.bounds, zoom = this.zoom;

        var lat1 = bounds.getSouthWest().lat(),
            lng1 = bounds.getSouthWest().lng(),
            lat2 = bounds.getNorthEast().lat(),
            lng2 = bounds.getNorthEast().lng();

        var tileB = TileSystem.getTileXY(lat1, lng1, zoom),
            tileT = TileSystem.getTileXY(lat2, lng2, zoom);

        var tileX, tileY, key, points = [], unvisitedQuadKeys = [];

        for (var i = 0, len1 = (tileT.x - tileB.x); i <= len1; ++i) {
            for (var j = 0, len2 = (tileB.y - tileT.y); j <= len2; ++j) {
                tileX = tileB.x + i;
                tileY = tileT.y + j;
                key = TileSystem.tileXYToQuadKey(tileX, tileY, zoom);

                if (this.cache.hasOwnProperty(key)) {
                    points.concat(this.cache[key]);
                } else {
                    unvisitedQuadKeys.push(key);
                    this.cache[key] = [];
                }
            }
        }

        if (points.length) {
            this.onPointReceived(points, true);
        }

        if (unvisitedQuadKeys.length) {
            this._send({
                lat1: lat1,
                lng1: lng1,
                lat2: lat2,
                lng2: lng2,
                quadKeys: unvisitedQuadKeys,
                zoom: zoom
            });
        }
    };

    return CachedRemoteManager;
});
