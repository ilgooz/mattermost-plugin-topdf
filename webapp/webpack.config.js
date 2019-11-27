var path = require('path');

module.exports = {
    entry: [
        './src/index.tsx',
    ],
    resolve: {
        modules: [
            'src',
            'node_modules',
        ],
        extensions: ['*', '.js', '.jsx', '.tsx'],
    },
    module: {
        rules: [
            {
              test: /\.(js|jsx|tsx)$/,
              exclude: /node_modules/,
              use: {
                  loader: 'babel-loader',
                  options: {
                      presets: ['@babel/preset-react',
                          [
                              "@babel/preset-env",
                              {
                                  "modules": "commonjs",
                                  "targets": {
                                      "node": "current"
                                  }
                              }
                          ]
                      ],
                  },
              },
          },
          {
            test: /\.(png|eot|tiff|svg|woff2|woff|ttf|gif|mp3|jpg)$/,
            use: [
                {
                    loader: 'url-loader',
                    options: {
                        name: 'files/[hash].[ext]',
                    },
                },
                {
                    loader: 'image-webpack-loader',
                    options: {},
                },
            ],
        },
        ],
    },
    externals: {
        react: 'React',
        'react-redux': 'ReactRedux',
        'prop-types': 'PropTypes'

    },
    output: {
        path: path.join(__dirname, '/dist'),
        publicPath: '/',
        filename: 'main.js',
    },
};