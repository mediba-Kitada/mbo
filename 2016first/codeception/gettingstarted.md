## [Getting Started](http://codeception.com/docs/02-GettingStarted)

> Let’s take a look at Codeception’s architecture. We assume that you already installed it, and bootstrapped your first test suites. Codeception has generated three of them: unit, functional, and acceptance. They are well described in the previous chapter. Inside your /tests folder you will have three config files and three directories with names corresponding to these suites. Suites are independent groups of tests with a common purpose.

それでは、Codeceptionのアーキテクチャを見ていきましょう。既に[インストール](http://codeception.com/install)が完了し、最初のテストをブートストラップしたとします。Codeceptionは、unit, Functional, acceptanceの３つのスイートを生成しているはずです。これらの詳細は、以降の章で見ていきます。```/test```ディレクトリには、３つの設定ファイルが格納されており、３つのスイート名に合致するディレクトリが作成されているはずです。スイートは、共通の目的別に独立しているテストのグループです。

### Actors

> One of the main concepts of Codeception is representation of tests as actions of a person.

codeceptionの特徴のひとつに、特定のユーザのアクションをテストとして実装することが挙げられます。

> We have a UnitTester, who executes functions and tests the code. We also have a FunctionalTester, a qualified tester, who tests the application as a whole, with knowledge of its internals. And an AcceptanceTester, a user that works with our application through an interface that we provide.

こう考えましょう。我々のチームには、コードをテストするテスターがいます。アプリケーション全体を内部の実装を意識したテストを実施出来る有能なテスターもいます。インターフェイスがワークすることを確認してくれる受入テスト担当者もいます。


> Actor classes are not written but generated from suite configuration. Methods of actor classes are generally taken from Codeception Modules. Each module provides predefined actions for different testing purposes, and they can be combined to fit the testing environment. Codeception tries to solve 90% of possible testing issues in its modules, so you don’t have reinvent the wheel. 

Actorクラス群は、コードとして記述されていませんが、codeceptionスイーツが設定によってジェネレートします。Actorクラスのメソッド群は、大抵の場合、Codeceptionモジュールにより利用されます。それらのモジュールは、異なる目的のテストのため再定義されたり、テスト環境の構築にために利用されたります。Codeceptionは、テストに関する問題の内、9割の解決に挑戦してますので、車輪を再開発する必要は無いわけです。

> We think that you can spend more time on writing tests and less on writing support code to make those tests run. By default AcceptanceTester relies on PhpBrowser module, which is set in tests/acceptance.suite.yml configuration file:

これまで以上にテストコードの開発に時間を割け、テストを実施するための余計なコードを書かずに済むようになるはずです。通常、受入テストモジュールは、```tests/acceptance.suite.yml```によってPhpBrowserモジュールと紐付けられています。

```yaml
class_name: AcceptanceTester
modules:
  enabled:
    - PhpBrowser:
      url: http://localhost/myapp/
    - \Helper\Acceptance
```

> In this configuration file you can enable/disable and reconfigure modules for your needs. When you change configuration, actor classes are rebuilt automatically. If Actor classes are not created or updated as you expect, try to generate them manually with build command:

この設定ファイルは、モジュールの利用可否を設定でき、必要に応じて再設定も出来ます。設定変更の際、Actorクラス群は、自動的にリビルドされます。もしアクタークラス群が意図通りに作成/更新されない場合は、以下のbuildコマンドで生成してみてください。

```bash
$ php codecept.phar build
```

### Writing a Sample Scenario

> By default tests are written as narrative scenarios. To make a PHP file a valid scenario, its name should have a Cept suffix.

通常、テストは物語風のシナリオとして記述していきます。Codeceptionで有効なシナリオをPHPとして実装する場合、Ceptをサフィックスを用います。

> Let’s say, we created a file tests/acceptance/SigninCept.php 

ここでは、```test/acceptance/SigninCept.php```を実装しましょう。

> We can do that by running the following command:

このシナリオを作成するため、以下のようなコマンドを実行します。

```bash
$ php codecept.phar generate:cept acceptance Signin
```

> A Scenario always starts with Actor class initialization. After that, writing a scenario is just like typing $I-> and choosing a proper action from the auto-completion list.

シナリオは必ずActorクラスの初期化を実行します。後は、$I->をタイプして、所定のアクションを自動補完のリストから選び、テストを書きまくります。

```php
<?php
$I = new AcceptanceTester($scenario);
?>
```

> Let’s log in to our website. We assume that we have a ‘login’ page where we are getting authenticated by providing a username and password. Then we are sent to a user page, where we see the text Hello, %username%. Let’s look at how this scenario is written in Codeception.

我々のウェブサイトにログインしてみましょう。ここでは、ユーザ名とパスワードによる認証機能をもつ'login'ページがあるとします。ログイン後は、'Hello, %username%' とページに表示されるとします。では、Codeceptionでのシナリオの書き方を見ていきましょう。

```php
<?php
$I = new AcceptanceTester($scenario);
$I->wantTo('log in as regular user');
$I->amOnPage('/login');
$I->fillField('Username','davert');
$I->fillField('Password','qwerty');
$I->click('Login');
$I->see('Hello, davert');
?>
```

> Before we execute this test, we should make sure that the website is running on a local web server. Let’s open the tests/acceptance.suite.yml file and replace the URL with the URL of your web application:

テストを実行する前に、ローカル環境でウェブサイトが稼働していることを確認しておきましょう。```tests/acceptance.suite.yml```を開き、URLを書き換えましょう。

```yaml
class_name: AcceptanceTester
modules:
  enabled:
    - PhpBrowser:
      url: 'http://myappurl.local'
    - \Helper\Acceptance
```

> After we configured the URL we can run this test with the run command:

設定したらrunコマンドでテストを実行してみます。

```bash
$ php codecept.phar run
```

> Here is the output we should see:

結果は、以下の様になります。

```bash
Acceptance Tests (1) -------------------------------
Trying log in as regular user (SigninCept.php)   Ok
----------------------------------------------------

Functional Tests (0) -------------------------------
----------------------------------------------------

Unit Tests (0) -------------------------------------
----------------------------------------------------

Time: 1 second, Memory: 21.00Mb

OK (1 test, 1 assertions)
```

> Let’s get a detailed output:

詳しく結果を見てみましょう。


```bash
$ php codecept.phar run acceptance --steps
```

> We should see a step-by-step report on the performed actions.

テストコードに対しステップバイステップで結果が確認出来るはずです。

```bash
Acceptance Tests (1) -------------------------------
Trying to log in as regular user (SigninCept.php)
Scenario:
* I am on page "/login"
* I fill field "Username" "davert"
* I fill field "Password" "qwerty"
* I click "Login"
* I see "Hello, davert"
  OK
----------------------------------------------------  

Time: 0 seconds, Memory: 21.00Mb

OK (1 test, 1 assertions)
```

> This simple test can be extended to a complete scenario of site usage. So by emulating the user’s actions you can test any of your websites.
> Give it a try!

このシンプルなテストは、実際のサイトで行われているような複雑なシナリオに拡張出来ます。ウェブサイトでのユーザアクションをエミュレートしてみましょう。
さあ、やってみましょう!!

### Bootstrap

> Each suite has its own bootstrap file. It’s located in the suite directory and is named _bootstrap.php. It will be executed before test suite. There is also a global bootstrap file located in the tests directory. It can be used to include additional files.

各スイートは、ブートストラップを持ち合わせています。ブートストラップは、```_bootstrap.php```の命名規則で```suite```ディレクトリに格納されています。```tests```ディレクトリには、グローバルなブートストラップファイルも格納されています。これは、追加ファイルをインクルードする際に利用されます。

### Cept, Cest and Test Formats

> Codeception supports three test formats. Beside the previously described scenario-based Cept format, Codeception can also execute PHPUnit test files for unit testing, and Cest format.

Codeceptionは、3種類のテストフォーマットをサポートしています。前述のシナリオベースのCeptフォーマットに加えてPHPUnite、Cestフォーマットによるテストをサポートしています。

> Cest combines scenario-driven test approach with OOP design. In case you want to group a few testing scenarios into one you should consider using Cest format. In the example below we are testing CRUD actions within a single file but with a several test (one per each operation):

Cestフォーマットは、シナリオ駆動テストアプローとをオブジェクト志向デザインを結合させたものです。とあるシナリオ群をとあるグループとしてまとめたい場合、Cestフォーマットを利用すべきです。以下の例では、CRUDアクションを単一ファイルに実装された複数のテストコードでテストするものです(one per each operation)。

```
<?php
class PageCrudCest
{
    function _before(AcceptanceTester $I)
    {
        // will be executed at the beginning of each test
        $I->amOnPage('/');
    }

    function createPage(AcceptanceTester $I)
    {
        // todo: write test
    }

    function viewPage(AcceptanceTester $I)
    {
        // todo: write test
    }

    function updatePage(AcceptanceTester $I)
    {
        // todo: write test
    }

    function deletePage(AcceptanceTester $I)
    {
        // todo: write test
    }
} 
?>
```

> Such Cest file can be created by running a generator:

このような```Cest```ファイルは、以下のようなコマンドで作成出来ます。

```bash
$ php codecept.phar generate:cest acceptance PageCrud
```

> Learn more about [Cest format](http://codeception.com/docs/07-AdvancedUsage#Cest-Classes) in Advanced Testing section.

詳しくは、発展的なテストの[Cestフォーマット](http://codeception.com/docs/07-AdvancedUsage#Cest-Classes)の欄を参考にしてください。

### Configuration

> Codeception has a global configuration in codeception.yml and a config for each suite. We also support .dist configuration files. If you have several developers in a project, put shared settings into codeception.dist.yml and personal settings into codeception.yml. The same goes for suite configs. For example, the unit.suite.yml will be merged with unit.suite.dist.yml.

Codeceptionは、```codeception.yml```でグローバルな設定情報を保持し、各スイートを設定します。```.dist```設定ファイルもサポートしています。複数の開発者がプロジェクトにアサインされている場合、```codeception.dist.yml```に共有すべき設定情報を記載し、```codeception.yml```に個人的な設定情報を記載します。スイートの設定も同様です。例えば、```unit.suite.yml```は、```unit.suite.dist.yml```にマージされます。

### Running Tests

> Tests can be started with the run command.

各テストは、```run```コマンドで実行出来ます。

```bash
$ php codecept.phar run
```

> With the first argument you can run tests from one suite.

第一引数で指定したスイートの各テストを実行出来ます。

```bash
$ php codecept.phar run acceptance
```

> To run exactly one test, add a second argument. Provide a local path to the test, from the suite directory.

特定のテストを実行するためには、第二引数で指定します。スイートのディレクトリからテストコードへのパスを渡します。

```bash
$ php codecept.phar run acceptance SigninCept.php
```

> Alternatively you can provide the full path to test file:

また、テストコードへのフルパスを渡すことも出来ます。

```bash
$ php codecept.phar run tests/acceptance/SigninCept.php
```

> You can execute one test from a test class (for Cest or Test formats)

Cestフォーマットの場合、テストクラスを指定しテストを実行することが出来ます。

```bash
$ php codecept.phar run tests/acceptance/SignInCest.php:anonymousLogin
```

> You can provide a directory path as well:
> This will execute all tests from the backend dir.

ディレクトリを渡すことも出来ます。  
ここでば```backend```ディレクトリに格納されている各テストを実行しています。

```bash
$ php codecept.phar run tests/acceptance/backend
```

> To execute a group of tests that are not stored in the same directory, you can organize them in groups.

同一ディレクトリに格納されていないテストをグループ単位で実行するため、[groups](http://codeception.com/docs/07-AdvancedUsage#Groups)を設定することが出来ます。

### Reports

> To generate JUnit XML output, you can provide the --xml option, and --html for HTML report.
> This command will run all tests for all suites, displaying the steps, and building HTML and XML reports. Reports will be stored in the tests/_output/ directory.

JUnit XMLフォーマットのレポートを生成するため、```--xml```オプションを付与することが出来ます。また、```--html```オプションでHTML形式のレポートを出力することが出来ます。  
このコマンドは、全スイートの全テストを実行し、ステップ毎に結果を表示し、XML及びHTMLのレポートを出力します。各レポートは、```tests/_output/```ディレクトリに格納されます。

```bash
$ php codecept.phar run --steps --xml --html
```

> To learn all available options, run the following command:

利用可能なオプションは、以下のコマンドで確認出来ます。

```bash
$ php codecept.phar help run
```

### Debugging

> To receive detailed output, tests can be executed with the --debug option. You may print any information inside a test using the codecept_debug function.

詳細な標準出力を受取るため、```--debug```オプションを付与しテストを実行することが出来ます。```codecept_debug```関数を用いてテスト内部の様々な情報を標準出力することが出来ます。

#### Generators

> There are plenty of useful Codeception commands:

一般的で便利コマンドCodeceptionコマンドは、以下の通りです。

> generate:cept suite filename - Generates a sample Cept scenario

- ```generate:cept suite filename``` - Ceptフォーマットのシナリオテストサンプルコードを生成します。

> generate:cest suite filename - Generates a sample Cest test

- ```generate:cest suite filename``` - Cestフォーマットのサンプルテストコードを生成します。

> generate:test suite filename - Generates a sample PHPUnit Test with Codeception hooks

- ```generate:test suite filename``` - Codeceptionフックを実装したPHPUnit形式のサンプルテストコードを生成します。

> generate:phpunit suite filename - Generates a classic PHPUnit Test

- ```generate:phpunit suite filename``` - 伝統的なPHPUnitテストコードを生成します。

> generate:suite suite actor - Generates a new suite with the given Actor class name

- ```generate:suite suite actor``` - 指定のActorクラス名で指定スイートを生成します。

> generate:scenarios suite - Generates text files containing scenarios from tests

- ```generate:scenarios suite - Generates``` 各テストからシナリオが記載されたテキストファイルを生成します。

> generate:helper filename - Generates a sample Helper File

- ```generate:helper filename``` - サンプルヘルパーを生成します。

> generate:pageobject suite filename - Generates a sample Page object

- ```generate:pageobject suite filename``` - ページオブジェクトのサンプルを生成します。

> generate:stepobject suite filename - Generates a sample Step object

- ```generate:stepobject suite filename``` - ステップオブジェクトのサンプルを生成します。

> generate:environment env - Generates a sample Environment configuration

- ```generate:environment env``` - 環境設定のサンプルを生成します。

> generate:groupobject group - Generates a sample Group Extension

- ```generate:groupobject group``` - グループエクステンションのサンプルを生成します。

### Conclusion

> We took a look into the Codeception structure. Most of the things you need were already generated by the bootstrap command. After you have reviewed the basic concepts and configurations, you can start writing your first scenario.

Codeceptionの構造について、見てきました。テストに必要なことは大抵ブートストラップによって生成されています。基本的なコンセプトと設定を俯瞰出来たら、最初のシナリオテストが書けるはずです。

