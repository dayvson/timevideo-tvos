var Presenter = {

    defaultPresenter: function(xml) {
        if(this.loadingIndicatorVisible) {
            navigationDocument.replaceDocument(xml, this.loadingIndicator);
            this.loadingIndicatorVisible = false;
        } else {
            navigationDocument.pushDocument(xml);
        }
    },
    menuBarItemPresenter: function(xml, ele) {
        var feature = ele.parentNode.getFeature("MenuBarDocument");

        if (feature) {
            var currentDoc = feature.getDocument(ele);
            if (!currentDoc) {
                feature.setDocument(xml, ele);
            }
        }
    },
    modalDialogPresenter: function(xml) {
        navigationDocument.presentModal(xml);
    },

    load: function(event) {
        var self = this,
        ele = event.target,
        videoURL = ele.getAttribute("videoURL")
        if(videoURL) {
            var player = new Player();
            var playlist = new Playlist();
            var mediaItem = new MediaItem("video", videoURL);

            player.playlist = playlist;
            player.playlist.push(mediaItem);
            player.present();
        }

        var self = this,
            ele = event.target,
            templateURL = ele.getAttribute("template"),
            presentation = ele.getAttribute("presentation");

        if (templateURL) {
            self.showLoadingIndicator(presentation);
            resourceLoader.loadResource(templateURL,
                function(resource) {
                    if (resource) {
                        var doc = self.makeDocument(resource);
                        doc.addEventListener("select", self.load.bind(self));
                        if (self[presentation] instanceof Function) {
                            self[presentation].call(self, doc, ele);
                        } else {
                            self.defaultPresenter.call(self, doc);
                        }
                    }
                }
            );
        }
    },

    makeDocument: function(resource) {
        if (!Presenter.parser) {
            Presenter.parser = new DOMParser();
        }

        var doc = Presenter.parser.parseFromString(resource, "application/xml");
        return doc;
    },
    showLoadingIndicator: function(presentation) {

        if (!this.loadingIndicator) {
            this.loadingIndicator = this.makeDocument(this.loadingTemplate);
        }
        if (!this.loadingIndicatorVisible && presentation != "modalDialogPresenter" && presentation != "menuBarItemPresenter") {
            navigationDocument.pushDocument(this.loadingIndicator);
            this.loadingIndicatorVisible = true;
        }
    },

    removeLoadingIndicator: function() {
        if (this.loadingIndicatorVisible) {
            navigationDocument.removeDocument(this.loadingIndicator);
            this.loadingIndicatorVisible = false;
        }
    },

    loadingTemplate: `<?xml version="1.0" encoding="UTF-8" ?>
        <document>
          <loadingTemplate>
            <activityIndicator>
              <text>Loading...</text>
            </activityIndicator>
          </loadingTemplate>
        </document>`
}