import { Store } from 'redux';
import { FileInfo } from 'mattermost-redux/types/files'

import OfficeDocsPreview, { isOfficeDocument } from './registers/office-docs-preview';

class TOPDFPlugin {
    initialize(registry, store: Store) {
        registry.registerFilePreviewComponent(
            ({ extension }: { extension: FileInfo } ) => isOfficeDocument(extension),
            OfficeDocsPreview(store),
        );
    }
}

window.registerPlugin('topdf', new TOPDFPlugin());
