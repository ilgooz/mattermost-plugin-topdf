import { Store } from 'redux';
import { FileInfo } from 'mattermost-redux/types/files'

import OfficeDocsPreview, { isOfficeDocument } from './registers/office-docs-preview';
import { setStore } from './delete/stores/redux_store';

class TOPDFPlugin {
    initialize(registry, store: Store) {
        setStore(store);
        
        registry.registerFilePreviewComponent(
            ({ extension }: { extension: FileInfo } ) => isOfficeDocument(extension),
            OfficeDocsPreview(store),
        );
    }
}

window.registerPlugin('topdf', new TOPDFPlugin());
