## mattermost-plugin-topdf
_TOPDF_ plugin automatically generates PDF previews of Microsoft Office files attached to posts for use with Mattermost. The plugin itself depends on [Gotenburg](https://github.com/thecodingmachine/gotenberg) for the heavy lifting, but seamlessly intercepts attachments to generate the previews without further user interaction.

## Installation
1. Install _Gotenberg_ server since this plugin depends on it to actually convert files to PDFs. Assuming you have already installed Docker, following command will download and install the latest version of _Gotenberg_ and set up a server with an open port at 4798:
  ```
    docker run -d -p 4798:3000 --env DEFAULT_WAIT_TIMEOUT=600 --env MAXIMUM_WAIT_TIMEOUT=600 thecodingmachine/gotenberg:6
  ```
  Make sure that you've set large timeouts for Gotenberg server like in the sample above since converting big files takes time.

  For more information about how to install & customize Gotenberg server, please follow the [docs](https://thecodingmachine.github.io/gotenberg/#install).

2. Go to the [releases page of this Github repository](https://github.com/ilgooz/mattermost-plugin-topdf/releases) and download the latest release for your Mattermost server.
   
3. In the Mattermost System Console under **System Console > Plugins > Plugin Management** upload the file to install the plugin. To learn more about how to upload a plugin, [see the documentation](https://docs.mattermost.com/administration/plugins.html#plugin-uploads).

4. Once _Gotenberg_ server is running, configure the plugin to make requests to your _Gotenberg_ instance. Go to **System Console > Plugins > TOPDF** and configure **Gotenberg's Full Address** to point at your _Gotenberg_ instance.  

5. Enabled the plugin at **System Console > Plugins > TOPDF** and ensure it starts with no errors by checking the server logs.

## Testing
To test your configuration is correct, post an Office file like _.docx_ to a channel and click on it to preview. If you see the content of document in a popup everything works fine.

# TODO
* `webapp/src/delete` should be deleted when `PDFPreview` component is accessible through Plugin API.
