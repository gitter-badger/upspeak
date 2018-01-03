'use strict'
const webpack = require('webpack')
const merge = require('webpack-merge')
const path = require('path')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const baseWebpackConfig = require('./webpack.base.conf')
const config = require('./buildConfig')

module.exports = merge(baseWebpackConfig, {
  module: {
    rules: [
      {
        test: /\.css$/,
        use: [
          'vue-style-loader',
          { loader: 'css-loader', options: { sourceMap: true } },
          { loader: 'postcss-loader', options: { sourceMap: true } }
        ]
      },
      {
        test: /\.less$/,
        use: [
          'vue-style-loader',
          { loader: 'css-loader', options: { sourceMap: true } },
          { loader: 'postcss-loader', options: { sourceMap: true } },
          { loader: 'less-loader', options: { sourceMap: true } }
        ]
      }
    ]
  },
  devtool: 'cheap-module-eval-source-map',
  plugins: [
    new webpack.DefinePlugin({
      'process.env': { NODE_ENV: '"development"' }
    }),
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NamedModulesPlugin(),
    new webpack.NoEmitOnErrorsPlugin(),
    new HtmlWebpackPlugin({
      filename: 'index.html',
      template: 'index.html',
      inject: true
    }),
    // copy custom static assets
    new CopyWebpackPlugin([
      {
        from: path.join(config.srcDir, 'static'),
        to: config.assetsDir,
        ignore: ['.*']
      }
    ])
  ]
})
