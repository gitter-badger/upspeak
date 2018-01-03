'use strict'
const path = require('path')
const webpack = require('webpack')
const merge = require('webpack-merge')
const CopyWebpackPlugin = require('copy-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const ExtractTextPlugin = require('extract-text-webpack-plugin')
const OptimizeCSSPlugin = require('optimize-css-assets-webpack-plugin')
const UglifyJsPlugin = require('uglifyjs-webpack-plugin')

const config = require('./buildConfig')
const baseWebpackConfig = require('./webpack.base.conf')

module.exports = merge(baseWebpackConfig, {
  module: {
    rules: [
      {
        test: /\.css$/,
        use: ExtractTextPlugin.extract({
          use: [
            { loader: 'css-loader', options: { sourceMap: false } },
            { loader: 'postcss-loader', options: { sourceMap: false } }
          ],
          fallback: 'vue-style-loader'
        })
      },
      {
        test: /\.less$/,
        use: ExtractTextPlugin.extract({
          use: [
            { loader: 'css-loader', options: { sourceMap: false } },
            { loader: 'postcss-loader', options: { sourceMap: false } },
            { loader: 'less-loader', options: { sourceMap: false } }
          ],
          fallback: 'vue-style-loader'
        })
      }
    ]
  },
  devtool: false,
  output: {
    path: config.outputDir,
    filename: path.join(config.assetsDir, 'js/[name].[chunkhash].js'),
    chunkFilename: path.join(config.assetsDir, 'js/[id].[chunkhash].js')
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env': { NODE_ENV: '"production"' }
    }),
    new UglifyJsPlugin({
      uglifyOptions: {
        compress: {
          warnings: false
        }
      },
      sourceMap: false,
      parallel: true
    }),
    // extract css into its own file
    new ExtractTextPlugin({
      filename: path.join(config.assetsDir, 'css/[name].[contenthash].css'),
      allChunks: true
    }),
    new OptimizeCSSPlugin({
      cssProcessorOptions: { safe: true }
    }),
    new HtmlWebpackPlugin({
      filename: config.indexPath,
      template: 'index.html',
      inject: true,
      minify: {
        removeComments: true,
        collapseWhitespace: true,
        removeAttributeQuotes: true
      },
      chunksSortMode: 'dependency'
    }),
    new webpack.HashedModuleIdsPlugin(),
    new webpack.optimize.ModuleConcatenationPlugin(),
    new webpack.optimize.CommonsChunkPlugin({
      name: 'vendor',
      children: false,
      minChunks (module) {
        return (module.context && module.context.includes('node_modules')) &&
          (module.resource && /\.js$/.test(module.resource))
      }
    }),
    new webpack.optimize.CommonsChunkPlugin({
      name: 'manifest',
      minChunks: Infinity
    }),
    new webpack.optimize.CommonsChunkPlugin({
      name: 'app',
      async: 'vendor-async',
      children: true,
      minChunks: 3
    }),
    new CopyWebpackPlugin([
      {
        from: path.join(config.srcDir, 'static'),
        to: config.assetsDir,
        ignore: ['.*']
      }
    ])
  ]
})
