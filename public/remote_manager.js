(function(root, factory) {
    if ("function" === typeof define && define.amd) {
        define("remote_manager", ["google_maps", "tile_system", "inherence"], factory);
    } else if ("object" === typeof module && module.exports) {
        module.exports = factory(
            require('google_maps'),
            require('tile_system'),
            require('inherence')
        );
    } else {
        root.RemoteManager = factory(
            root.google,
            root.TileSystem,
            root.Inherence
        );
    }
})(this, function(google, TileSystem, Inherence) {
    var RemoteManager = function(map, wsPath, onPointRecieved) {
        this.map = map;
        this.ws = new WebSocket(wsPath);
        this.onPointRecieved = onPointRecieved;
    };

    RemoteManager.prototype.run = function() {
        this.map.addListener('idle', this._onIdle.bind(this));
        this.ws.onmessage = this._onMessage.bind(this);
    };

    RemoteManager.prototype._onIdle = function() {
        this._updateState();
        this._fetchPoints();
    };

    RemoteManager.prototype._fetchPoints = function() {
        var bounds = this.bounds, zoom = this.zoom;

        var lat1 = bounds.getSouthWest().lat()
            lng1 = bounds.getSouthWest().lng(),
            lat2 = bounds.getNorthEast().lat(),
            lng2 = bounds.getNorthEast().lng();

        var tileB = TileSystem.getTileXY(lat1, lng1, zoom),
            tileT = TileSystem.getTileXY(lat2, lng2, zoom);

        this._send({
            lat1: lat1,
            lng1: lng1,
            lat2: lat2,
            lng2: lng2,
            tileBounds: [tileB.x, tileB.y, tileT.x, tileT.y],
            zoom: zoom
        });
    };

    RemoteManager.prototype._send = function(toSend) {
        console.log(toSend);
        this.ws.send(JSON.stringify(toSend))
    };

    RemoteManager.prototype._onMessage = function(event) {
        var data = JSON.parse(event.data);
        this.onPointRecieved(data);
    };

    RemoteManager.prototype._updateState = function() {
        this.zoom = this.map.getZoom();
        this.bounds = this.map.getBounds();
    };

    var CachedRemoteManager = (function() {
        Inherence(CachedRemoteManager, RemoteManager);

        function CachedRemoteManager() {
            CachedRemoteManager.__super__.constructor.apply(this, arguments);
            this.cache = {};
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

    CachedRemoteManager.prototype._fetchPoints = function() {
        var bounds = this.bounds, zoom = this.zoom;

        var lat1 = bounds.getSouthWest().lat()
            lng1 = bounds.getSouthWest().lng(),
            lat2 = bounds.getNorthEast().lat(),
            lng2 = bounds.getNorthEast().lng();

        var tileB = TileSystem.getTileXY(lat1, lng1, zoom),
            tileT = TileSystem.getTileXY(lat2, lng2, zoom);

        var tileX, tileY, key, points = [];

        for (var i = 0, len1 = (tileT.x - tileB.x); i < len1; ++i) {
            for (var j = 0, len2 = (tileT.y - tileB.y); i < len2; ++j) {
                tileX = tileB.x + i;
                tileY = tileB.y + j;
                key = tileX + '_' + tileY;

                if (this.cache.hasOwnProperty(key)) {
                    points.concat(this.cache[key]);
                } else {

                }
            }
        }

        this._send({
            lat1: lat1,
            lng1: lng1,
            lat2: lat2,
            lng2: lng2,
            tileBounds: [tileB.x, tileB.y, tileT.x, tileT.y],
            zoom: zoom
        });
    };


    return CachedRemoteManager;
});
