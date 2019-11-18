// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

type Props = {
    text: React.ReactNode;
}

export default class LoadingSpinner extends React.PureComponent<Props> {
    public static defaultProps: Props = {
        text: null,
    }

    public render() {
        return (
            <span
                id='loadingSpinner'
                className={'LoadingSpinner' + (this.props.text ? ' with-text' : '')}
            >
                <span
                    className='fa fa-spinner fa-fw fa-pulse spinner'
                    title='Loading Icon'
                />
                {this.props.text}
            </span>
        );
    }
}
