require 'health-data-standards'
require 'open-uri'
require 'highline/import'
require 'mongoid'

Mongoid.configure do |config|
  config.sessions = { default: { hosts: [ "localhost:27017" ], database: 'fhir' }}
end

namespace :bundle do
  desc %{ Download measure/test deck bundle.
      options
      nlm_user    - the nlm username to authenticate to the server - will prompt is not supplied
      nlm_passwd  - the nlm password for authenticating to the server - will prompt if not supplied
      version     - the version of the bundle to download. This will default to the version

     example usage:
      rake bundle:download nlm_name=username nlm_passwd=password version=2.1.0-latest
  }
  task :download do
    nlm_user = ENV["nlm_user"]
    nlm_passwd = ENV["nlm_pass"]
    measures_dir = File.join(Dir.pwd, "bundles")

    while nlm_user.nil? || nlm_user == ""
      nlm_user = ask("NLM Username?: "){ |q| q.readline = true }
    end

    while nlm_passwd.nil? || nlm_passwd == ""
      nlm_passwd = ask("NLM Password?: "){ |q| q.echo = false
                 q.readline = true }
    end

    bundle_version = ENV["version"] || "2016"
    @bundle_name = "bundle-#{bundle_version}.zip"

    puts "Downloading and saving #{@bundle_name} to #{measures_dir}"
    # Pull down the list of bundles and download the version we're looking for
    bundle_uri = "https://cypressdemo.healthit.gov/measure_bundles/#{@bundle_name}"
    bundle = nil

    tries = 0
    max_tries = 10
    last_error = nil
    while bundle.nil? && tries < max_tries do
      tries = tries + 1
      begin
        bundle = open(bundle_uri, :proxy => ENV["http_proxy"],:http_basic_authentication=>[nlm_user, nlm_passwd] )
            rescue OpenURI::HTTPError => oe
        last_error = oe
        if oe.message == "401 Unauthorized"
          puts "Please check your credentials and try again"
          break
        end
            rescue => e
        last_error = e
        sleep 0.5
      end
    end

    if bundle.nil?
       puts "An error occured while downloading the bundle"
      raise last_error if last_error
    end
    # Save the bundle to the measures directory
    FileUtils.mkdir_p measures_dir
    FileUtils.mv(bundle.path, File.join(measures_dir, @bundle_name))

  end

  desc %{ Download and install the measure/test deck bundle.  This is essientally delegating to the bundle_download and bundle:import tasks
    options
    nlm_user    - the nlm username to authenticate to the server - will prompt is not supplied
    nlm_passwd  - the nlm password for authenticating to the server - will prompt if not supplied
    version     - the version of the bundle to download. This will default to the version
    delete_existing - delete any existing bundles with the same version and reinstall - default is false - will cause error if same version already exists
    update_measures - update any existing measures with the same hqmf_id to those contained in this bundle.
          Will only work for bundle versions greater than that of the installed version - default is false
    type -  type of measures to be installed from bundle. A bundle may have measures of different types such as ep or eh.  This will constrain the types installed, defautl is all types
   example usage:
    rake budnle:download_and_install nlm_name=username nlm_passwd=password version=2.1.0-latest  type=ep
  }
  task :download_and_install => [:download] do
    de = ENV['delete_existing'] || false
    um = ENV['update_measures'] || false
    puts "Importing bundle #{@bundle_name} delete_existing: #{de}  update_measures: #{um} type: #{ENV['type'] || 'ALL'}"
    task("bundle:import").invoke("bundles/#{@bundle_name}",de, um , ENV['type'], 'true')
  end

  desc 'Import a quality bundle into the database.'
  task :import, [:bundle_path,  :delete_existing,  :update_measures, :type, :create_indexes, :exclude_results] do |task, args|
    raise "The path to the measures zip file must be specified" unless args.bundle_path
    options = {:delete_existing => (args.delete_existing == "true"),
               :type => args.type,
               :update_measures => (args.update_measures == "true"),
               :exclude_results => (args.exclude_results == "true")
              }

    bundle = File.open(args.bundle_path)
    importer = HealthDataStandards::Import::Bundle::Importer
    bundle_contents = importer.import(bundle, options)

    counts = {measures: bundle_contents.measures.count,
              records: bundle_contents.records.count,
              extensions: bundle_contents[:extensions].count,
              value_sets: bundle_contents.value_sets.count}

    if (args.create_indexes != 'false')
      ::Rails.application.eager_load! if defined? Rails
      ::Mongoid::Tasks::Database.create_indexes
    end

    puts "Successfully imported bundle at: #{args.bundle_path}"
    puts "\t Imported into environment: #{Rails.env.upcase}" if defined? Rails
    puts "\t Loaded #{args.type || 'all'} measures"
    puts "\t Sub-Measures Loaded: #{counts[:measures]}"
    puts "\t Test Patients Loaded: #{counts[:records]}"
    puts "\t Extensions Loaded: #{counts[:extensions]}"
    puts "\t Value Sets Loaded: #{counts[:value_sets]}"
  end

end
