// Webpack build configuration parameters

const path = require('path')

const srcDir = path.resolve(__dirname, '../')
const outputDir = path.resolve(srcDir, '../bin/web')
const isProduction = process.env.NODE_ENV === 'production'
// Path where index.html will be produced
const indexPath = path.join(outputDir, 'index.html')
// Name of assets subdirectory within the build path
const assetsDir = 'static'
const publicPath = '/'

module.exports = {
  srcDir,
  outputDir,
  isProduction,
  indexPath,
  assetsDir,
  publicPath
}
