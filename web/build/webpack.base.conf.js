'use strict'
const path = require('path')
const config = require('./buildConfig')

module.exports = {
  context: config.srcDir,
  entry: {
    app: './main.js'
  },
  output: {
    path: config.outputDir,
    filename: '[name].js',
    publicPath: config.publicPath
  },
  resolve: {
    extensions: ['.js', '.vue', '.json'],
    alias: {
      '@': config.srcDir
    }
  },
  module: {
    rules: [
      {
        test: /\.(js|vue)$/,
        loader: 'eslint-loader',
        enforce: 'pre',
        include: [config.srcDir],
        exclude: [path.join(config.srcDir, 'node_modules')],
        options: {
          emitWarning: false
        }
      },
      {
        test: /\.vue$/,
        loader: 'vue-loader',
        options: {
          loaders: [
            'vue-style-loader',
            {
              loader: 'css-loader',
              options: { sourceMap: !config.isProduction }
            },
            {
              loader: 'less-loader',
              options: { sourceMap: !config.isProduction }
            }
          ],
          cssSourceMap: true,
          cacheBusting: true,
          transformToRequire: {
            video: ['src', 'poster'],
            source: 'src',
            img: 'src',
            image: 'xlink:href'
          }
        }
      },
      {
        test: /\.js$/,
        loader: 'babel-loader',
        include: [config.srcDir, path.resolve('node_modules/webpack-dev-server/client')],
        exclude: [path.resolve(config.srcDir, 'node_modules/')]
      },
      {
        test: /\.(png|jpe?g|gif|svg)(\?.*)?$/,
        loader: 'url-loader',
        options: {
          limit: 10000,
          name: path.join(config.assetsDir, 'img/[name].[hash:7].[ext]')
        }
      },
      {
        test: /\.(mp4|webm|ogg|mp3|wav|flac|aac)(\?.*)?$/,
        loader: 'url-loader',
        options: {
          limit: 10000,
          name: path.join(config.assetsDir, 'media/[name].[hash:7].[ext]')
        }
      },
      {
        test: /\.(woff2?|eot|ttf|otf)(\?.*)?$/,
        loader: 'url-loader',
        options: {
          limit: 10000,
          name: path.join(config.assetsDir, 'fonts/[name].[hash:7].[ext]')
        }
      }
    ]
  },
  node: {
    setImmediate: false,
    dgram: 'empty',
    fs: 'empty',
    net: 'empty',
    tls: 'empty',
    child_process: 'empty'
  }
}
